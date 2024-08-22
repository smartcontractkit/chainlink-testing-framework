package havoc

import (
	"context"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ChaosEntity is an interface that defines common behaviors for chaos management entities.
type ChaosEntity interface {
	// Create initializes and submits the chaos object to Kubernetes.
	Create(ctx context.Context)
	// Delete removes the chaos object from Kubernetes.
	Delete(ctx context.Context) error
	// Registers a listener to receive updates about the chaos object's lifecycle.
	AddListener(listener ChaosListener)

	GetObject() client.Object
	GetChaosName() string
	GetChaosDescription() string
	GetChaosDuration() (time.Duration, error)
	GetChaosSpec() interface{}
	GetStartTime() time.Time
	GetEndTime() time.Time
	GetExpectedEndTime() (time.Time, error)
}
