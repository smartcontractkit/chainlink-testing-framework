package client

import (
	"context"
	"fmt"
	"time"

	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/chaos-mesh/chaos-mesh/api/v1alpha1"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/smartcontractkit/chainlink-env/config"
)

// Chaos is controller that manages Chaosmesh CRD instances to run experiments
type Chaos struct {
	Client         *K8sClient
	ResourceByName map[string]string
	Namespace      string
}

type ChaosState struct {
	ChaosDetails v1alpha1.ChaosStatus `json:"status"`
}

// NewChaos creates controller to run and stop chaos experiments
func NewChaos(client *K8sClient, namespace string) *Chaos {
	return &Chaos{
		Client:         client,
		ResourceByName: make(map[string]string),
		Namespace:      namespace,
	}
}

// Run runs experiment and saves its ID
func (c *Chaos) Run(app cdk8s.App, id string, resource string) (string, error) {
	log.Info().Msg("Applying chaos experiment")
	config.JSIIGlobalMu.Lock()
	manifest := *app.SynthYaml()
	config.JSIIGlobalMu.Unlock()
	log.Trace().Str("Raw", manifest).Msg("Manifest")
	c.ResourceByName[id] = resource
	if err := c.Client.Apply(context.Background(), manifest, c.Namespace); err != nil {
		return id, err
	}
	if err := c.checkForPodsExistence(app); err != nil {
		return id, err
	}
	err := c.waitForChaosStatus(id, v1alpha1.ConditionAllInjected, time.Minute)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (c *Chaos) waitForChaosStatus(id string, condition v1alpha1.ChaosConditionType, timeout time.Duration) error {
	var result ChaosState
	log.Info().Msgf("waiting for chaos experiment state %s", condition)
	if timeout < time.Minute {
		log.Info().Msg("timeout is less than 1 minute, setting to 1 minute")
		timeout = time.Minute
	}
	return wait.PollImmediate(2*time.Second, timeout, func() (bool, error) {
		data, err := c.Client.ClientSet.
			RESTClient().
			Get().
			RequestURI(fmt.Sprintf("/apis/chaos-mesh.org/v1alpha1/namespaces/%s/%s/%s", c.Namespace, c.ResourceByName[id], id)).
			Do(context.Background()).
			Raw()
		if err == nil {
			err = json.Unmarshal(data, &result)
			if err != nil {
				return false, err
			}
			for _, c := range result.ChaosDetails.Conditions {
				if c.Type == condition && c.Status == v1.ConditionTrue {
					return true, err
				}
			}
		}
		return false, nil
	})
}

func (c *Chaos) WaitForAllRecovered(id string, timeout time.Duration) error {
	return c.waitForChaosStatus(id, v1alpha1.ConditionAllRecovered, timeout)
}

// Stop removes a chaos experiment
func (c *Chaos) Stop(id string) error {
	defer delete(c.ResourceByName, id)
	return c.Client.DeleteResource(c.Namespace, c.ResourceByName[id], id)
}

func (c *Chaos) checkForPodsExistence(app cdk8s.App) error {
	charts := app.Charts()
	var selectors []string
	for _, chart := range *charts {
		json := chart.ToJson()
		for _, j := range *json {
			m := j.(map[string]interface{})
			fmt.Println(m)
			kind := m["kind"].(string)
			if kind == "PodChaos" || kind == "NetworkChaos" {
				selectors = append(selectors, getLabelSelectors(m["spec"].(map[string]interface{})))
			}
			if kind == "NetworkChaos" {
				target := m["spec"].(map[string]interface{})["target"].(map[string]interface{})
				selectors = append(selectors, getLabelSelectors(target))
			}
		}
	}
	for _, selector := range selectors {
		podList, err := c.Client.ListPods(c.Namespace, selector)
		if err != nil {
			return err
		}
		if podList == nil || len(podList.Items) == 0 {
			return fmt.Errorf("no pods found for selector %s", selector)
		}
		log.Info().
			Int("podsCount", len(podList.Items)).
			Str("selector", selector).
			Msgf("found pods for chaos experiment")
	}
	return nil
}

func getLabelSelectors(spec map[string]interface{}) string {
	if spec == nil {
		return ""
	}
	s := spec["selector"].(map[string]interface{})
	if s == nil {
		return ""
	}
	m := s["labelSelectors"].(map[string]interface{})
	selector := ""
	for key, value := range m {
		if selector == "" {
			selector = fmt.Sprintf("%s=%s", key, value)
		} else {
			selector = fmt.Sprintf("%s, %s=%s", selector, key, value)
		}
	}
	return selector
}
