package cleaner

import (
	"context"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

// TimeoutPolicy struct for timeout policy for various e2e tests
type TimeoutPolicy struct {
	Namespace v1.Namespace
	client    *kubernetes.Clientset
	Applied   bool
}

// IsApplied checks if policy is already applied
func (p *TimeoutPolicy) IsApplied() bool {
	return p.Applied
}

// Apply removes namespace if timeout is reached
func (p *TimeoutPolicy) Apply() error {
	if p.Applied {
		return nil
	}
	elapsed := time.Since(p.Namespace.CreationTimestamp.Time)
	timeout, ok := p.Namespace.Labels["timeout"]
	if !ok {
		log.Warn().Str("Namespace", p.Namespace.Name).Msg("No timeout field found")
		p.Applied = true
		return nil
	}
	timeoutDur, err := time.ParseDuration(timeout)
	if err != nil {
		return err
	}
	canApply := elapsed > timeoutDur
	log.Info().
		Str("Namespace", p.Namespace.Name).
		Str("Elapsed", elapsed.String()).
		Bool("CanApply", canApply).
		Msg("Watching test namespace")
	if !canApply {
		return nil
	}
	p.Applied = true
	if err := p.client.CoreV1().Namespaces().Delete(context.Background(), p.Namespace.Name, v12.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}
