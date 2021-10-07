package cleaner

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/environment"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"time"
)

// NamespacePolicy arbitrary set of rules that can by applied to test namespaces
type NamespacePolicy interface {
	Apply() error
	IsApplied() bool
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// Config cleaner service config
type Config struct {
	PollInterval time.Duration
}

// Cleaner cleaner service struct with a set of policies to apply
type Cleaner struct {
	cfg      *Config
	client   *kubernetes.Clientset
	policies map[string]NamespacePolicy
}

// NewCleaner creates new cleaner service
func NewCleaner(client *kubernetes.Clientset, cfg *Config) *Cleaner {
	return &Cleaner{
		cfg:      cfg,
		client:   client,
		policies: map[string]NamespacePolicy{},
	}
}

// updatePolicies gets all test namespaces marked by particular labels and create policies for a new one
func (c *Cleaner) updatePolicies() error {
	ns, err := c.client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{
		LabelSelector: environment.BasicTestNamespaceSelector,
		FieldSelector: environment.NamespaceActivePhaseSelector,
	})
	if err != nil {
		return err
	}
	for _, n := range ns.Items {
		typ, ok := n.ObjectMeta.Labels["policy"]
		if !ok {
			log.Error().Str("Namespace", n.Name).Msg("No type label found on")
			continue
		}
		switch typ {
		case "timeout":
			c.policies[n.Name] = &TimeoutPolicy{
				Namespace: n,
				client:    c.client,
			}
		}
	}
	return nil
}

// Run runs cleaner loop
func (c *Cleaner) Run() error {
	for {
		if err := c.updatePolicies(); err != nil {
			return err
		}
		for ns, policy := range c.policies {
			if policy.IsApplied() {
				delete(c.policies, ns)
			}
			if err := policy.Apply(); err != nil {
				log.Err(err).Send()
				return err
			}
		}
		log.Debug().Int("Total", len(c.policies)).Msg("Total policies")
		time.Sleep(c.cfg.PollInterval)
	}
}
