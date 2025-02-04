package chaos

import (
	"context"
	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/havoc"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
)

func defaultLogger() zerolog.Logger {
	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.TraceLevel)
}

type NodeLatenciesConfig struct {
	Namespace               string
	Description             string
	Latency                 time.Duration
	InjectionDuration       time.Duration
	FromLabelKey            string
	FromLabelValues         []string
	ToLabelKey              string
	ToLabelValues           []string
	ExperimentTotalDuration time.Duration
	ExperimentCreateDelay   time.Duration
}

func networkDelay(client client.Client, l zerolog.Logger, cfg NodeLatenciesConfig) (*havoc.Schedule, error) {
	return havoc.NewSchedule(havoc.ScheduleOpts{
		Description: cfg.Description,
		DelayCreate: cfg.ExperimentCreateDelay,
		Duration:    cfg.ExperimentTotalDuration,
		Logger:      &l,
		Object: &v1alpha1.Schedule{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Schedule",
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "latencies",
				Namespace: cfg.Namespace,
			},
			Spec: v1alpha1.ScheduleSpec{
				Schedule:          "@every 1m",
				ConcurrencyPolicy: v1alpha1.ForbidConcurrent,
				Type:              v1alpha1.ScheduleTypeNetworkChaos,
				HistoryLimit:      2,
				ScheduleItem: v1alpha1.ScheduleItem{
					EmbedChaos: v1alpha1.EmbedChaos{
						NetworkChaos: &v1alpha1.NetworkChaosSpec{
							Action: v1alpha1.DelayAction,
							PodSelector: v1alpha1.PodSelector{
								Mode: v1alpha1.AllMode,
								Selector: v1alpha1.PodSelectorSpec{
									GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
										Namespaces: []string{cfg.Namespace},
										ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
											{
												Operator: "In",
												Key:      cfg.FromLabelKey,
												Values:   cfg.FromLabelValues,
											},
										},
									},
								},
							},
							Direction: v1alpha1.From,
							Target: &v1alpha1.PodSelector{
								Selector: v1alpha1.PodSelectorSpec{
									GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
										Namespaces: []string{cfg.Namespace},
										ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
											{
												Operator: "In",
												Key:      cfg.ToLabelKey,
												Values:   cfg.ToLabelValues,
											},
										},
									},
								},
								Mode: v1alpha1.AllMode,
							},
							TcParameter: v1alpha1.TcParameter{
								Delay: &v1alpha1.DelaySpec{
									Latency:     cfg.Latency.String(),
									Correlation: "100",
									Jitter:      "0ms",
								},
							},
							Duration: ptr.To[string]((cfg.InjectionDuration).String()),
						},
					},
				},
			},
		},
		Client: client,
	})
}

type NodeRebootsConfig struct {
	Namespace               string
	Description             string
	LabelKey                string
	LabelValues             []string
	InjectionDuration       time.Duration
	ExperimentTotalDuration time.Duration
	ExperimentCreateDelay   time.Duration
}

func podFail(client client.Client, l zerolog.Logger, cfg NodeRebootsConfig) (*havoc.Schedule, error) {
	return havoc.NewSchedule(havoc.ScheduleOpts{
		Description: cfg.Description,
		DelayCreate: cfg.ExperimentCreateDelay,
		Duration:    cfg.ExperimentTotalDuration,
		Logger:      &l,
		Object: &v1alpha1.Schedule{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Schedule",
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fail",
				Namespace: cfg.Namespace,
			},
			Spec: v1alpha1.ScheduleSpec{
				Schedule:          "@every 1m",
				ConcurrencyPolicy: v1alpha1.ForbidConcurrent,
				Type:              v1alpha1.ScheduleTypePodChaos,
				HistoryLimit:      2,
				ScheduleItem: v1alpha1.ScheduleItem{
					EmbedChaos: v1alpha1.EmbedChaos{
						PodChaos: &v1alpha1.PodChaosSpec{
							Action: v1alpha1.PodFailureAction,
							ContainerSelector: v1alpha1.ContainerSelector{
								PodSelector: v1alpha1.PodSelector{
									Mode: v1alpha1.AllMode,
									Selector: v1alpha1.PodSelectorSpec{
										GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
											Namespaces: []string{cfg.Namespace},
											ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
												{
													Operator: "In",
													Key:      cfg.LabelKey,
													Values:   cfg.LabelValues,
												},
											},
										},
									},
								},
							},
							Duration: ptr.To[string](cfg.InjectionDuration.String()),
						},
					},
				},
			},
		},
		Client: client,
	})
}

type NodeCPUStressConfig struct {
	Namespace               string
	Description             string
	Cores                   int
	CoreLoadPercentage      int // 0-100
	LabelKey                string
	LabelValues             []string
	InjectionDuration       time.Duration
	ExperimentTotalDuration time.Duration
	ExperimentCreateDelay   time.Duration
}

func podStressCPU(client client.Client, l zerolog.Logger, cfg NodeCPUStressConfig) (*havoc.Schedule, error) {
	return havoc.NewSchedule(havoc.ScheduleOpts{
		Description: cfg.Description,
		DelayCreate: cfg.ExperimentCreateDelay,
		Duration:    cfg.ExperimentTotalDuration,
		Logger:      &l,
		Object: &v1alpha1.Schedule{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Schedule",
				APIVersion: "chaos-mesh.org/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "stress",
				Namespace: cfg.Namespace,
			},
			Spec: v1alpha1.ScheduleSpec{
				Schedule:          "@every 1m",
				ConcurrencyPolicy: v1alpha1.ForbidConcurrent,
				Type:              v1alpha1.ScheduleTypeStressChaos,
				HistoryLimit:      2,
				ScheduleItem: v1alpha1.ScheduleItem{
					EmbedChaos: v1alpha1.EmbedChaos{
						StressChaos: &v1alpha1.StressChaosSpec{
							ContainerSelector: v1alpha1.ContainerSelector{
								PodSelector: v1alpha1.PodSelector{
									Mode: v1alpha1.AllMode,
									Selector: v1alpha1.PodSelectorSpec{
										GenericSelectorSpec: v1alpha1.GenericSelectorSpec{
											Namespaces: []string{cfg.Namespace},
											ExpressionSelectors: v1alpha1.LabelSelectorRequirements{
												{
													Operator: "In",
													Key:      cfg.LabelKey,
													Values:   cfg.LabelValues,
												},
											},
										},
									},
								},
							},
							Stressors: &v1alpha1.Stressors{
								CPUStressor: &v1alpha1.CPUStressor{
									Stressor: v1alpha1.Stressor{
										Workers: cfg.Cores,
									},
									Load: ptr.To[int](cfg.CoreLoadPercentage),
								},
							},
							Duration: ptr.To[string](cfg.InjectionDuration.String()),
						},
					},
				},
			},
		},
		Client: client,
	})
}

func TestK8sChaos(t *testing.T) {
	l := defaultLogger()
	c, err := havoc.NewChaosMeshClient()
	require.NoError(t, err)

	// example
	namespace := "crib-radek-ccip-v2"
	experimentDuration := 3 * time.Minute

	rebootsChaos, err := podFail(c, l, NodeRebootsConfig{
		Namespace:               namespace,
		Description:             "reboot nodes",
		LabelKey:                "app.kubernetes.io/instance",
		LabelValues:             []string{"ccip-1"},
		ExperimentTotalDuration: experimentDuration,
	})
	require.NoError(t, err)
	latenciesChaos, err := networkDelay(c, l, NodeLatenciesConfig{
		Namespace:               namespace,
		Description:             "network issues",
		Latency:                 5000 * time.Millisecond,
		InjectionDuration:       40 * time.Second,
		FromLabelKey:            "app.kubernetes.io/instance",
		FromLabelValues:         []string{"ccip-1"},
		ToLabelKey:              "app.kubernetes.io/instance",
		ToLabelValues:           []string{"ccip-2", "ccip-3", "ccip-4"},
		ExperimentTotalDuration: 3 * time.Minute,
		ExperimentCreateDelay:   3 * time.Minute,
	})
	require.NoError(t, err)
	blockchainLatency, err := networkDelay(c, l, NodeLatenciesConfig{
		Namespace:               namespace,
		Description:             "blockchain nodes network issues",
		Latency:                 5000 * time.Millisecond,
		InjectionDuration:       40 * time.Second,
		FromLabelKey:            "app.kubernetes.io/instance",
		FromLabelValues:         []string{"ccip-1", "ccip-2", "ccip-3"},
		ToLabelKey:              "instance",
		ToLabelValues:           []string{"geth-1337"},
		ExperimentTotalDuration: 3 * time.Minute,
		//ExperimentCreateDelay:   3 * time.Minute,
	})
	require.NoError(t, err)
	podCPUStress, err := podStressCPU(c, l, NodeCPUStressConfig{
		Namespace:               namespace,
		Description:             "CPU hog",
		LabelKey:                "app.kubernetes.io/instance",
		LabelValues:             []string{"ccip-1"},
		Cores:                   1,
		CoreLoadPercentage:      90,
		ExperimentTotalDuration: experimentDuration,
	})
	require.NoError(t, err)

	_ = rebootsChaos
	_ = latenciesChaos
	_ = blockchainLatency
	_ = podCPUStress

	chaosList := []havoc.ChaosEntity{
		rebootsChaos,
		latenciesChaos,
		blockchainLatency,
		podCPUStress,
	}

	for _, chaos := range chaosList {
		chaos.AddListener(havoc.NewChaosLogger(l))
		exists, err := havoc.ChaosObjectExists(chaos.GetObject(), c)
		require.NoError(t, err)
		require.False(t, exists, "chaos object already exists: %s. Delete it before starting the test", chaos.GetChaosName())
		chaos.Create(context.Background())
	}
	t.Cleanup(func() {
		for _, chaos := range chaosList {
			chaos.Delete(context.Background())
		}
	})

	// Start WASP load test here, apply average load profile that you expect in production!
	// Configure timeouts and validate all the test cases until the test ends
	time.Sleep(3*time.Minute + 5*time.Second)
}
