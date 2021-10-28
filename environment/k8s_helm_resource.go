package environment

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/kube"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/portforward"
)

const (
	// HelmInstallTimeout timeout for installing a helm chart
	HelmInstallTimeout = 200 * time.Second
	// ReleasePrefix the default prefix
	ReleasePrefix = "release"
	// DefaultK8sConfigPath the default path for kube
	DefaultK8sConfigPath = ".kube/config"
)

// SetValuesHelmFunc interface for setting values in a helm chart
type SetValuesHelmFunc func(resource *HelmChart) error

// PodForwardedInfo data to port forward the pods
type PodForwardedInfo struct {
	PodIP          string
	ForwardedPorts []portforward.ForwardedPort
	PodName        string
}

// HelmChart celoextended helm chart data
type HelmChart struct {
	id                string
	chartPath         string
	releaseName       string
	actionConfig      *action.Configuration
	env               *K8sEnvironment
	network           *config.NetworkConfig
	SetValuesHelmFunc SetValuesHelmFunc
	// Deployment properties
	pods         []PodForwardedInfo
	values       map[string]interface{}
	stopChannels []chan struct{}
}

// Teardown tears down the helm release
func (k *HelmChart) Teardown() error {
	// closing forwarded ports
	for _, stopChan := range k.stopChannels {
		stopChan <- struct{}{}
	}
	log.Debug().Str("Release", k.releaseName).Msg("Uninstalling Helm release")
	if _, err := action.NewUninstall(k.actionConfig).Run(k.releaseName); err != nil {
		return err
	}
	return nil
}

// ID returns the helm chart id
func (k *HelmChart) ID() string {
	return k.id
}

// SetValue sets the specified value in the chart
func (k *HelmChart) SetValue(key string, val interface{}) {
	k.values[key] = val
}

// GetConfig gets the helms environment config
func (k *HelmChart) GetConfig() *config.Config {
	return k.env.config
}

// Values returns the helm charts values
func (k *HelmChart) Values() map[string]interface{} {
	return k.values
}

// SetEnvironment sets the environment
func (k *HelmChart) SetEnvironment(environment *K8sEnvironment) error {
	k.env = environment
	return nil
}

// Environment gets environment
func (k *HelmChart) Environment() *K8sEnvironment {
	return k.env
}

func (k *HelmChart) forwardAllPodsPorts() error {
	k8sPods := k.env.k8sClient.CoreV1().Pods(k.env.namespace.Name)
	pods, err := k8sPods.List(context.Background(), metaV1.ListOptions{
		LabelSelector: k.releaseSelector(),
	})
	if err != nil {
		return err
	}
	for _, p := range pods.Items {
		ports, err := forwardPodPorts(&p, k.env.k8sConfig, k.env.namespace.Name, k.stopChannels)
		if err != nil {
			return fmt.Errorf("unable to forward ports: %v", err)
		}
		k.pods = append(k.pods, PodForwardedInfo{
			PodIP:          p.Status.PodIP,
			ForwardedPorts: ports,
			PodName:        p.Name,
		})
		log.Info().Str("Manifest ID", k.id).Interface("Ports", ports).Msg("Forwarded ports")
	}
	return nil
}

// WaitUntilHealthy waits until the helm release is healthy
func (k *HelmChart) WaitUntilHealthy() error {
	// using helm Wait option before, not need to wait for pods to be deployed there
	if err := k.forwardAllPodsPorts(); err != nil {
		return err
	}
	if k.values == nil {
		k.values = make(map[string]interface{})
	}
	if k.SetValuesHelmFunc != nil {
		if err := k.SetValuesHelmFunc(k); err != nil {
			return err
		}
	}
	return nil
}

func (k *HelmChart) releaseSelector() string {
	return fmt.Sprintf("%s=%s", ReleasePrefix, k.releaseName)
}

// ServiceDetails gets the details of the released service
func (k *HelmChart) ServiceDetails() ([]*ServiceDetails, error) {
	var serviceDetails []*ServiceDetails
	for _, pod := range k.pods {
		for _, port := range pod.ForwardedPorts {
			remoteURL, err := url.Parse(fmt.Sprintf("http://%s:%d", pod.PodIP, port.Remote))
			if err != nil {
				return serviceDetails, err
			}
			localURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port.Local))
			if err != nil {
				return serviceDetails, err
			}
			serviceDetails = append(serviceDetails, &ServiceDetails{
				RemoteURL: remoteURL,
				LocalURL:  localURL,
			})
		}
	}
	return serviceDetails, nil
}

// Deploy deploys the helm charts
func (k *HelmChart) Deploy(_ map[string]interface{}) error {
	log.Info().Str("Path", k.chartPath).
		Str("Release", k.releaseName).
		Str("Namespace", k.env.namespace.Name).
		Msg("Installing Helm chart")
	chart, err := loader.Load(k.chartPath)
	if err != nil {
		return err
	}

	chart.Values, err = chartutil.CoalesceValues(chart, k.values)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	k.actionConfig = &action.Configuration{}

	// TODO: So, this is annoying, and not really all that important, I SHOULD be able to just use our K8sConfig function
	// and pass that in as our config, but K8s has like 10 different config types, all of which don't talk to each other,
	// and this wants an interface, instead of the rest config that we use everywhere else. Creating such an interface is
	// also a huge hassle and... well anyway, if you've got some time to burn to make this more sensical, I hope you like
	// digging into K8s code with sparse to no docs.
	kubeConfigPath := filepath.Join(homeDir, DefaultK8sConfigPath)
	if len(os.Getenv("KUBECONFIG")) > 0 {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
	if err := k.actionConfig.Init(
		kube.GetConfig(kubeConfigPath, "", k.env.namespace.Name),
		k.env.namespace.Name,
		os.Getenv("HELM_DRIVER"),
		func(format string, v ...interface{}) {
			log.Debug().Str("LogType", "Helm").Msg(fmt.Sprintf(format, v...))
		}); err != nil {
		return err
	}

	install := action.NewInstall(k.actionConfig)
	install.Namespace = k.env.namespace.Name
	install.ReleaseName = k.releaseName
	install.Timeout = HelmInstallTimeout
	// blocks until all pods are healthy
	install.Wait = true
	_, err = install.Run(chart, nil)
	if err != nil {
		return err
	}
	log.Info().
		Str("Namespace", k.env.namespace.Name).
		Str("Release", k.releaseName).
		Str("Chart", k.chartPath).
		Msg("Succesfully installed helm chart")
	return nil
}
