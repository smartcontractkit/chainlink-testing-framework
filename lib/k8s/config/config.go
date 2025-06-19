package config

import (
	"os"
	"sync"

	"dario.cat/mergo"
)

const (
	EnvVarPrefix = "TEST_"

	E2ETestEnvVarPrefix = "E2E_TEST_"

	EnvVarSkipManifestUpdate            = "SKIP_MANIFEST_UPDATE"
	EnvVarSkipManifestUpdateDescription = "Skip updating manifest when connecting to the namespace"
	EnvVarSkipManifestUpdateExample     = "false"

	EnvVarKeepEnvironments            = "KEEP_ENVIRONMENTS"
	EnvVarKeepEnvironmentsDescription = "Should we keep environments on test completion"
	EnvVarKeepEnvironmentsExample     = "NEVER|ALWAYS|ON_FAILURE"

	EnvVarNamespace            = "ENV_NAMESPACE"
	EnvVarNamespaceDescription = "Namespace name to connect to"
	EnvVarNamespaceExample     = "chainlink-test-epic"

	EnvVarJobImage            = "ENV_JOB_IMAGE"
	EnvVarJobImageDescription = "Image to run as a job in k8s"
	EnvVarJobImageExample     = "795953128386.dkr.ecr.us-west-2.amazonaws.com/core-integration-tests:v1.0"

	EnvVarInsideK8s            = "ENV_INSIDE_K8S"
	EnvVarInsideK8sDescription = "Internal variable to turn forwarding strategy off inside k8s, do not use"
	EnvVarInsideK8sExample     = ""

	// deprecated (use TOML config instead to pass the image)
	EnvVarCLImage            = "CHAINLINK_IMAGE"
	EnvVarCLImageDescription = "Chainlink image repository"
	EnvVarCLImageExample     = "public.ecr.aws/chainlink/chainlink"

	// deprecated (use TOML config instead to pass the version)
	EnvVarCLTag            = "CHAINLINK_VERSION"
	EnvVarCLTagDescription = "Chainlink image tag"
	EnvVarCLTagExample     = "1.9.0"

	EnvVarUser            = "CHAINLINK_ENV_USER"
	EnvVarUserDescription = "Owner of an environment"
	EnvVarUserExample     = "Satoshi"

	EnvVarTeam            = "CHAINLINK_USER_TEAM"
	EnvVarTeamDescription = "Team to, which owner of the environment belongs to"
	EnvVarTeamExample     = "BIX, CCIP, BCM"

	EnvVarCLCommitSha            = "CHAINLINK_COMMIT_SHA"
	EnvVarCLCommitShaDescription = "The sha of the commit that you're running tests on. Mostly used for CI"
	EnvVarCLCommitShaExample     = "${{ github.sha }}"

	EnvVarTestTrigger            = "TEST_TRIGGERED_BY"
	EnvVarTestTriggerDescription = "How the test was triggered, either manual or CI."
	EnvVarTestTriggerExample     = "CI"

	EnvVarLogLevel            = "TEST_LOG_LEVEL"
	EnvVarLogLevelDescription = "Environment logging level"
	EnvVarLogLevelExample     = "info | debug | trace"

	EnvVarDBURL            = "DATABASE_URL"
	EnvVarDBURLDescription = "DATABASE_URL needed for component test. This is only necessary if testhelper methods are imported from core"
	EnvVarDBURLExample     = "postgresql://postgres:node@localhost:5432/chainlink_test?sslmode=disable"

	EnvVarSlackKey            = "SLACK_API_KEY"
	EnvVarSlackKeyDescription = "The OAuth Slack API key to report tests results with"
	EnvVarSlackKeyExample     = "xoxb-example-key"

	EnvVarSlackChannel            = "SLACK_CHANNEL"
	EnvVarSlackChannelDescription = "The Slack code for the channel you want to send the notification to"
	EnvVarSlackChannelExample     = "C000000000"

	EnvVarSlackUser            = "SLACK_USER"
	EnvVarSlackUserDescription = "The Slack code for the user you want to notify"
	EnvVarSlackUserExample     = "U000000000"

	EnvVarToleration                 = "K8S_TOLERATION"
	EnvVarTolerationsUserDescription = "Node roles to tolerate"
	EnvVarTolerationsExample         = "foundations"

	EnvVarNodeSelector                = "K8S_NODE_SELECTOR"
	EnvVarNodeSelectorUserDescription = "Node role to deploy to"
	EnvVarNodeSelectorExample         = "foundations"

	EnvVarDetachRunner                = "DETACH_RUNNER"
	EnvVarDetachRunnerUserDescription = "Should we detach the remote runner after starting a test using it"
	EnvVarDetachRunnerExample         = "true"

	EnvVarRemoteRunnerCpu                = "RR_CPU"
	EnvVarRemoteRunnerCpuUserDescription = "The cpu limit and req for the remote runner"
	EnvVarRemoteRunnerCpuExample         = "1000m"

	EnvVarRemoteRunnerMem                = "RR_MEM"
	EnvVarRemoteRunnerMemUserDescription = "The mem limit and req for the remote runner"
	EnvVarRemoteRunnerMemExample         = "1024Mi"

	EnvVarInternalDockerRepo            = "INTERNAL_DOCKER_REPO"
	EnvVarInternalDockerRepoDescription = "Use internal docker repository for some images"
	EnvVarInternalDockerRepoExample     = "public.ecr.aws"

	EnvVarLocalCharts                = "LOCAL_CHARTS"
	EnvVarLocalChartsUserDescription = "Use local charts from the CTF repository directly"
	EnvVarLocalChartsExample         = "true"

	EnvBase64ConfigOverride             = "BASE64_CONFIG_OVERRIDE"
	EnvBase64ConfigOverriderDescription = "Base64-encoded TOML config (should contain at least chainlink image and version)"
	EnvBase64ConfigOverrideExample      = "W0NoYWlubGlua0ltYWdlXQppbWFnZT0icHVibGljLmVjci5hd3MvY2hhaW5saW5rL2NoYWlubGluayIKdmVyc2lvbj0iMi43LjEtYXV0b21hdGlvbi0yMDIzMTEyNyIKCltBdXRvbWF0aW9uXQpbQXV0b21hdGlvbi5HZW5lcmFsXQpkdXJhdGlvbj0yMDAK"

	EnvSethLogLevel            = "SETH_LOG_LEVEL"
	EnvSethLogLevelDescription = "Specifies the log level used by Seth"
	EnvSethLogLevelExample     = "info"
)

var (
	JSIIGlobalMu = &sync.Mutex{}
)

func MustMerge(targetVars interface{}, codeVars interface{}) {
	if err := mergo.Merge(targetVars, codeVars, mergo.WithOverride); err != nil {
		panic(err)
	}
}

func MustEnvOverrideVersion(target interface{}) {
	image := os.Getenv(EnvVarCLImage)
	tag := os.Getenv(EnvVarCLTag)
	if image != "" && tag != "" {
		if err := mergo.Merge(target, map[string]interface{}{
			"chainlink": map[string]interface{}{
				"image": map[string]interface{}{
					"image":   image,
					"version": tag,
				},
			},
		}, mergo.WithOverride); err != nil {
			panic(err)
		}
	}
}
