// sd-merge: discover Docker containers by label, fetch each container's /discovery JSON,
// enrich targets with labels, merge + dedupe, and write a single Prometheus file_sd JSON.
//
// Env (with defaults)
//   LABEL_MATCH=framework=ctf         # docker ps --filter "label=$LABEL_MATCH"
//   DISCOVERY_PATH=/discovery         # path inside the container that returns target groups
//   DISCOVERY_PORT=6688               # discovery port
//   DISCOVERY_SCHEME=http             # http or https
//   NETWORK_NAME=                     # prefer IP from this Docker network (optional)
//   OUT=/out/merged.json              # output path for file_sd
//   SLEEP=15                          # seconds or Go duration (e.g. 15s, 1m)
//   REQUEST_TIMEOUT=5                 # seconds or Go duration
//   REWRITE_TO_IP=0                   # 1 = replace target host with container IP, keep target port

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

// Prometheus file_sd/http_sd target group.
type TargetGroup struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type Config struct {
	LabelMatch    string
	DefaultPath   string
	DefaultPort   string
	DefaultScheme string
	PreferNetwork string
	OutPath       string
	Sleep         time.Duration
	ReqTimeout    time.Duration
	RewriteToIP   bool
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseDur(s string, def time.Duration) time.Duration {
	if s == "" {
		return def
	}
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	if n, err := strconv.Atoi(s); err == nil && n > 0 {
		return time.Duration(n) * time.Second
	}
	return def
}

func loadConfig() Config {
	return Config{
		LabelMatch:    getenv("LABEL_MATCH", "framework=ctf"),
		DefaultPath:   getenv("DISCOVERY_PATH", "/discovery"),
		DefaultPort:   getenv("DISCOVERY_PORT", "6688"),
		DefaultScheme: getenv("DISCOVERY_SCHEME", "http"),
		PreferNetwork: getenv("NETWORK_NAME", ""),
		OutPath:       getenv("OUT", "/out/merged.json"),
		Sleep:         parseDur(getenv("SLEEP", "15"), 15*time.Second),
		ReqTimeout:    parseDur(getenv("REQUEST_TIMEOUT", "5"), 5*time.Second),
		RewriteToIP:   getenv("REWRITE_TO_IP", "0") == "1",
	}
}

func main() {
	log.SetFlags(0)
	cfg := loadConfig()
	log.Printf("[sd-merge] start label=%q out=%s every=%s", cfg.LabelMatch, cfg.OutPath, cfg.Sleep)

	if err := os.MkdirAll(filepath.Dir(cfg.OutPath), 0o755); err != nil {
		log.Fatalf("mkdir out: %v", err)
	}
	_ = atomicWriteJSON(cfg.OutPath, []TargetGroup{}) // ensure file exists

	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		log.Fatalf("docker client: %v", err)
	}
	defer cli.Close()

	httpClient := &http.Client{Timeout: cfg.ReqTimeout}
	ctx := context.Background()

	ticker := time.NewTicker(cfg.Sleep)
	defer ticker.Stop()

	for {
		if err := runOnce(ctx, cli, httpClient, cfg); err != nil {
			log.Printf("[sd-merge] cycle error: %v", err)
		}
		<-ticker.C
	}
}

// One full discovery+merge cycle.
func runOnce(ctx context.Context, cli *docker.Client, hc *http.Client, cfg Config) error {
	ids, err := listContainerIDs(ctx, cli, cfg.LabelMatch)
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		if err := atomicWriteJSON(cfg.OutPath, []TargetGroup{}); err != nil {
			return err
		}
		log.Printf("[sd-merge] no matching containers -> wrote empty list")
		return nil
	}

	var all []TargetGroup

	for _, id := range ids {
		inspect, err := cli.ContainerInspect(ctx, id)
		if err != nil {
			log.Printf("[sd-merge] skip %.12s: inspect failed: %v", id, err)
			continue
		}
		ip := pickIP(inspect, cfg.PreferNetwork)
		if ip == "" {
			log.Printf("[sd-merge] skip %.12s: no IP", id)
			continue
		}
		name := strings.TrimPrefix(inspect.Name, "/")

		// Per-container overrides, with sane defaults.
		path := first(inspect.Config.Labels["prom_sd_path"], cfg.DefaultPath)
		port := first(inspect.Config.Labels["prom_sd_port"], cfg.DefaultPort)
		scheme := first(inspect.Config.Labels["prom_sd_scheme"], cfg.DefaultScheme)
		url := fmt.Sprintf("%s://%s:%s%s", scheme, ip, port, path)

		// Fetch discovery JSON (array of target groups).
		tgs, err := fetchDiscovery(hc, url)
		if err != nil {
			log.Printf("[sd-merge] %s fetch failed: %v", url, err)
			continue
		}

		// Enrich labels and (optionally) rewrite hosts to container IP.
		for i := range tgs {
			if tgs[i].Labels == nil {
				tgs[i].Labels = map[string]string{}
			}
			if mp := tgs[i].Labels["__metrics_path__"]; mp != "" {
				tgs[i].Labels["scrape_path"] = mp
			} else {
				tgs[i].Labels["scrape_path"] = ""
			}
			tgs[i].Labels["container_name"] = name

			if cfg.RewriteToIP {
				for j, tgt := range tgs[i].Targets {
					if p := targetPort(tgt); p != "" {
						tgs[i].Targets[j] = ip + ":" + p
					}
				}
			}
		}

		all = append(all, tgs...)
		log.Printf("[sd-merge] ok %s -> %d groups", url, len(tgs))
	}

	merged := mergeAndDedupe(all)
	if err := atomicWriteJSON(cfg.OutPath, merged); err != nil {
		return err
	}
	if fi, err := os.Stat(cfg.OutPath); err == nil {
		log.Printf("[sd-merge] wrote %d groups (%d bytes) -> %s", len(merged), fi.Size(), cfg.OutPath)
	}
	return nil
}

// listContainerIDs returns IDs of running containers matching a label filter.
func listContainerIDs(ctx context.Context, cli *docker.Client, label string) ([]string, error) {
	label = strings.TrimSpace(label)
	if label == "" {
		return nil, fmt.Errorf("LABEL_MATCH cannot be empty")
	}
	f := filters.NewArgs()
	f.Add("label", label)
	cs, err := cli.ContainerList(ctx, container.ListOptions{Filters: f})
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(cs))
	for _, c := range cs {
		ids = append(ids, c.ID)
	}
	return ids, nil
}

// pickIP chooses the preferred network IP, or the first available.
func pickIP(info types.ContainerJSON, prefer string) string {
	if info.NetworkSettings == nil || info.NetworkSettings.Networks == nil {
		return ""
	}
	if prefer != "" {
		if es, ok := info.NetworkSettings.Networks[prefer]; ok && es.IPAddress != "" {
			return es.IPAddress
		}
	}
	for _, es := range info.NetworkSettings.Networks {
		if es.IPAddress != "" {
			return es.IPAddress
		}
	}
	return ""
}

// fetchDiscovery GETs the discovery URL and decodes a JSON array of target groups.
// Any non-200 or non-array payload results in an empty list (not an error).
func fetchDiscovery(hc *http.Client, url string) ([]TargetGroup, error) {
	resp, err := hc.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
	}
	var out []TargetGroup
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return []TargetGroup{}, nil
	}
	return out, nil
}

func targetPort(hostport string) string {
	if i := strings.LastIndex(hostport, ":"); i >= 0 && i+1 < len(hostport) {
		return hostport[i+1:]
	}
	return ""
}

func first(v, def string) string {
	if v != "" {
		return v
	}
	return def
}

// mergeAndDedupe: group by identical label sets and union targets (sorted).
func mergeAndDedupe(in []TargetGroup) []TargetGroup {
	byKey := make(map[string]map[string]struct{}) // labelsKey -> targets set
	lblOf := make(map[string]map[string]string)

	for _, g := range in {
		k := labelsKey(g.Labels)
		if byKey[k] == nil {
			byKey[k] = make(map[string]struct{})
			lblOf[k] = g.Labels
		}
		for _, t := range g.Targets {
			byKey[k][t] = struct{}{}
		}
	}

	keys := make([]string, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out := make([]TargetGroup, 0, len(keys))
	for _, k := range keys {
		ts := make([]string, 0, len(byKey[k]))
		for t := range byKey[k] {
			ts = append(ts, t)
		}
		sort.Strings(ts)
		out = append(out, TargetGroup{Targets: ts, Labels: lblOf[k]})
	}
	return out
}

// labelsKey builds a canonical string for a label map (sorted key=value lines).
func labelsKey(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(m[k])
		b.WriteString("\n")
	}
	return b.String()
}

// atomicWriteJSON writes JSON to a temp file and renames it into place.
func atomicWriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
