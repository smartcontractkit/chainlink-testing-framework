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
// datasource-managed definitions from every discovered source. wantAll fetches
// every ds rule; otherwise only the named rules are requested. Ruler-returned ds
// rules are dropped — only discovery knows their datasource UID.
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

// dsFilterNames extracts the server-side rule_name[] filters. A key: or uid:
// form forces a bulk fetch; otherwise both the whole input and its last
// /-separated segment are sent (a ds rule's name can itself contain "/", while
// Group/Title names it by the trailing segment).
func dsFilterNames(names []string) (filters []string, fetchAll bool) {
	seen := make(map[string]bool)
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			filters = append(filters, s)
		}
	}
	for _, raw := range names {
		n := strings.TrimSpace(raw)
		if n == "" {
			continue
		}
		if strings.HasPrefix(n, "key:") || strings.HasPrefix(n, "uid:") {
			return nil, true
		}
		add(n)
		if i := strings.LastIndex(n, "/"); i != -1 {
			add(n[i+1:])
		}
	}
	return filters, false
}
