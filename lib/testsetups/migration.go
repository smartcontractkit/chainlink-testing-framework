package testsetups

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/chainlink"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/pkg/helm/ethereum"
)

type FromVersionSpec struct {
	Image string
	Tag   string
}

type ToVersionSpec struct {
	Image string
	Tag   string
}

type DBMigrationSpec struct {
	FromSpec          FromVersionSpec
	ToSpec            ToVersionSpec
	KeepConnection    bool
	RemoveOnInterrupt bool
}

// DBMigration returns an environment with DB migrated from FromVersionSpec to ToVersionSpec
func DBMigration(spec *DBMigrationSpec) (*environment.Environment, error) {
	e := environment.New(nil).
		AddHelm(ethereum.New(nil)).
		AddHelm(chainlink.New(0, map[string]interface{}{
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   spec.FromSpec.Image,
					"version": spec.FromSpec.Tag,
				},
			},
			"db": map[string]interface{}{
				"stateful": true,
				"capacity": "1Gi",
			},
		}))
	err := e.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to setup initial deployment for version: %s:%s err: %w", spec.FromSpec.Image, spec.FromSpec.Tag, err)
	}
	e.Cfg.KeepConnection = spec.KeepConnection
	e.Cfg.RemoveOnInterrupt = spec.RemoveOnInterrupt
	e.Cfg.UpdateWaitInterval = 10 * time.Second
	env, err := e.
		ReplaceHelm("chainlink-0", chainlink.New(0, map[string]interface{}{
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   spec.ToSpec.Image,
					"version": spec.ToSpec.Tag,
				},
			},
			"db": map[string]interface{}{
				"stateful": true,
				"capacity": "1Gi",
			},
		}))
	if err != nil {
		return nil, fmt.Errorf("failed to replace helm chart for version: %s:%s err: %w", spec.ToSpec.Image, spec.ToSpec.Tag, err)
	}
	err = env.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to migrate to version: %s:%s err: %w", spec.ToSpec.Image, spec.ToSpec.Tag, err)
	}
	return e, nil
}
