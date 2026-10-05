package gate

import (
	"context"
	"fmt"
	"strings"
)

// ListAllDefinitions is the CLI's entry point for `list`: every Grafana-managed
// and datasource-managed definition, discovered from scratch.
func ListAllDefinitions(ctx context.Context, src Source) ([]Definition, error) {
	return loadDefinitions(ctx, src, nil, true)
}

// loadDefinitions reads Grafana-managed definitions from the ruler and
// datasource-managed definitions from every discovered rule source. wantAll
// fetches every ds rule (label selection, list); otherwise only the named rules
// are requested, one request per source. Ruler-returned datasource-managed rules
// are dropped: discovery is the authority for them, since only it knows the
// datasource UID.
func loadDefinitions(ctx context.Context, src Source, names []string, wantAll bool) ([]Definition, error) {
	grafana, err := src.GrafanaDefinitions(ctx)
	if err != nil {
		return nil, err
	}
	sources, err := src.DiscoverRuleSources(ctx)
	if err != nil {
		return nil, err
	}

	defs := make([]Definition, 0, len(grafana))
	for _, d := range grafana {
		if d.Kind == KindDatasourceManaged {
			continue
		}
		defs = append(defs, d)
	}

	filter, fetchAll := dsFilterNames(names)
	if wantAll {
		fetchAll = true
	}
	for _, rs := range sources {
		var dsDefs []Definition
		if fetchAll {
			dsDefs, err = src.DatasourceDefinitions(ctx, rs, nil)
		} else {
			dsDefs, err = src.DatasourceDefinitions(ctx, rs, filter)
		}
		if err != nil {
			return nil, fmt.Errorf("datasource %q: %w", rs.Name, err)
		}
		defs = append(defs, dsDefs...)
	}
	return defs, nil
}

// dsFilterNames extracts the server-side rule_name[] filters from the raw alert
// names. A key: or uid: form names no title, so the whole source must be
// fetched; otherwise the last /-separated segment is the title.
func dsFilterNames(names []string) (titles []string, fetchAll bool) {
	for _, raw := range names {
		n := strings.TrimSpace(raw)
		if n == "" {
			continue
		}
		if strings.HasPrefix(n, "key:") || strings.HasPrefix(n, "uid:") {
			return nil, true
		}
		parts := strings.Split(n, "/")
		titles = append(titles, parts[len(parts)-1])
	}
	return titles, false
}
