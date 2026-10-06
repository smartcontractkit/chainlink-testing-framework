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

	fetch := planDatasourceFetch(names, wantAll)
	if !fetch.Skip {
		for _, rs := range sources {
			var dsDefs []Definition
			if fetch.All {
				dsDefs, err = src.DatasourceDefinitions(ctx, rs, nil)
			} else {
				dsDefs, err = src.DatasourceDefinitions(ctx, rs, fetch.Names)
			}
			if err != nil {
				return nil, fmt.Errorf("datasource %q: %w", rs.Name, err)
			}
			defs = append(defs, dsDefs...)
		}
	}
	// Duplicate keys are deliberately NOT rejected here: loading is an
	// inventory (`list`), and a collision only matters when both rules are
	// actually selected. resolveAlertSet enforces that on the selected set.
	return defs, nil
}

// dsFetchPlan is how loadDefinitions should query the datasource rule sources.
type dsFetchPlan struct {
	// Skip is set when no selector can name a datasource rule (every selector
	// is uid:, which only a Grafana rule carries), so the sources need not be
	// read at all.
	Skip bool
	// All fetches every rule of each source: a key: selector names a ds rule and
	// the API has no key filter, and wantAll asks for everything.
	All bool
	// Names are the rule_name[] filters for a name selection.
	Names []string
}

// planDatasourceFetch decides what to ask the datasource sources for. A uid:
// selector can only name a Grafana rule, so it never triggers a ds read; a key:
// selector can name a ds rule and forces a bulk read.
func planDatasourceFetch(names []string, wantAll bool) dsFetchPlan {
	if wantAll {
		return dsFetchPlan{All: true}
	}
	var plan dsFetchPlan
	seen := make(map[string]bool)
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			plan.Names = append(plan.Names, s)
		}
	}
	for _, raw := range names {
		n := strings.TrimSpace(raw)
		if n == "" {
			continue
		}
		switch {
		case strings.HasPrefix(n, "key:"):
			return dsFetchPlan{All: true}
		case strings.HasPrefix(n, "uid:"):
			continue
		}
		// Both the whole input and its last /-separated segment: a ds rule's
		// name can itself contain "/", while Group/Title names it by the
		// trailing segment.
		add(n)
		if i := strings.LastIndex(n, "/"); i != -1 {
			add(n[i+1:])
		}
	}
	if len(plan.Names) == 0 {
		return dsFetchPlan{Skip: true}
	}
	return plan
}
