package actions

import (
	"github.com/onsi/ginkgo"
	"github.com/smartcontractkit/helmenv/environment"
	"github.com/smartcontractkit/integrations-framework/client"
)

// Keep Environments options
const (
	KeepEnvironmentsNever  = "never"
	KeepEnvironmentsOnFail = "onfail"
	KeepEnvironmentsAlways = "always"
)

// TeardownSuite tears down networks/clients and environment
func TeardownSuite(env *environment.Environment, nets *client.Networks) error {
	if ginkgo.CurrentSpecReport().Failed() {
		// nolint
		if err := env.Artifacts.DumpTestResult(ginkgo.CurrentGinkgoTestDescription().FullTestText, "chainlink"); err != nil {
			return err
		}
	}
	if nets != nil {
		if err := nets.Teardown(); err != nil {
			return err
		}
	}
	if err := env.Teardown(); err != nil {
		return err
	}

	//switch strings.ToLower(config.KeepEnvironments) {
	//case KeepEnvironmentsNever:
	//	env.TearDown()
	//case KeepEnvironmentsOnFail:
	//	if !ginkgo.CurrentGinkgoTestDescription().Failed {
	//		env.TearDown()
	//	} else {
	//		log.Info().Str("Namespace", env.ID()).Msg("Kept environment due to test failure")
	//	}
	//case KeepEnvironmentsAlways:
	//	log.Info().Str("Namespace", env.ID()).Msg("Kept environment")
	//	return
	//default:
	//	env.TearDown()
	//}
	return nil
}
