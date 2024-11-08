package havoc

import (
	"context"
	"os"
	"time"

	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/rs/zerolog"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func ExampleNewChaos() {
	testLogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	client := fake.NewFakeClient()
	podChaos := &v1alpha1.PodChaos{ /* PodChaos spec */ }
	chaos, err := NewChaos(ChaosOpts{
		Object:      podChaos,
		Description: "Pod failure example",
		DelayCreate: 5 * time.Second,
		Client:      client,
		Logger:      &testLogger,
	})
	if err != nil {
		panic(err)
	}
	logger := NewChaosLogger(testLogger)
	annotator := NewSingleLineGrafanaAnnotator(
		"http://grafana-instance.com",
		"grafana-access-token",
		"dashboard-uid",
		Logger,
	)
	chaos.AddListener(logger)
	chaos.AddListener(annotator)

	chaos.Create(context.Background())
	// Output:
}
