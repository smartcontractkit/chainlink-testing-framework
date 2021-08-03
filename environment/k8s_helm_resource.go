package environment

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/integrations-framework/config"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	HelmInstallTimeout   = 200 * time.Second
	ReleasePrefix        = "release"
	DefaultK8sConfigPath = ".kube/config"
)

type SetValuesHelmFunc func(resource *HelmChart) error

type PodForwardedInfo struct {
	PodIP          string
	ForwardedPorts []portforward.ForwardedPort
}

// HelmChart common helm chart data
type HelmChart struct {
	id                string
	chartPath         string
	releaseName       string
	actionConfig      *action.Configuration
	k8sClient         *kubernetes.Clientset
	config            *config.Config
	k8sConfig         *rest.Config
	network           *config.NetworkConfig
	namespace         *coreV1.Namespace
	SetValuesHelmFunc SetValuesHelmFunc
	// Deployment properties
	pods         []PodForwardedInfo
	values       map[string]interface{}
	stopChannels []chan struct{}
}

func (k *HelmChart) Teardown() error {
	// closing forwarded ports
	for _, stopChan := range k.stopChannels {
		stopChan <- struct{}{}
	}
	log.Debug().Str("Release", k.releaseName).Msg("Uninstalling Рelm release")
	if _, err := action.NewUninstall(k.actionConfig).Run(k.releaseName); err != nil {
		return err
	}
	return nil
}

func (k *HelmChart) ID() string {
	return k.id
}

func (k *HelmChart) SetValue(key string, val interface{}) {
	k.values[key] = val
}

func (k *HelmChart) GetConfig() *config.Config {
	return k.config
}

func (k *HelmChart) Values() map[string]interface{} {
	return k.values
}

func (k *HelmChart) SetEnvironment(
	k8sClient *kubernetes.Clientset,
	k8sConfig *rest.Config,
	config *config.Config,
	network *config.NetworkConfig,
	namespace *coreV1.Namespace,
) error {
	k.k8sClient = k8sClient
	k.k8sConfig = k8sConfig
	k.config = config
	k.network = network
	k.namespace = namespace
	return nil
}

func (k *HelmChart) forwardAllPodsPorts() error {
	k8sPods := k.k8sClient.CoreV1().Pods(k.namespace.Name)
	pods, err := k8sPods.List(context.Background(), metaV1.ListOptions{
		LabelSelector: k.releaseSelector(),
	})
	if err != nil {
		return err
	}
	for _, p := range pods.Items {
		ports, err := forwardPodPorts(&p, k.k8sConfig, k.namespace.Name, k.stopChannels)
		if err != nil {
			return fmt.Errorf("unable to forward ports: %v", err)
		}
		k.pods = append(k.pods, PodForwardedInfo{
			PodIP:          p.Status.PodIP,
			ForwardedPorts: ports,
		})
		log.Info().Str("Manifest ID", k.id).Interface("Ports", ports).Msg("Forwarded ports")
	}
	return nil
}

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

func (k *HelmChart) Deploy(_ map[string]interface{}) error {
	log.Info().Str("Path", k.chartPath).Str("Namespace", k.namespace.Name).Msg("Installing helm chart")
	chart, err := loader.Load(k.chartPath)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	k.actionConfig = &action.Configuration{}

	if err := k.actionConfig.Init(
		kube.GetConfig(filepath.Join(homeDir, DefaultK8sConfigPath), "", k.namespace.Name),
		k.namespace.Name,
		os.Getenv("HELM_DRIVER"),
		func(format string, v ...interface{}) {
			log.Debug().Str("LogType", "Helm").Msg(fmt.Sprintf(format, v...))
		}); err != nil {
		return err
	}

	install := action.NewInstall(k.actionConfig)
	install.Namespace = k.namespace.Name
	install.ReleaseName = k.releaseName
	install.Timeout = HelmInstallTimeout
	// blocks until all pods are healthy
	install.Wait = true
	_, err = install.Run(chart, nil)
	if err != nil {
		return err
	}
	log.Info().
		Str("Namespace", k.namespace.Name).
		Str("Release", k.releaseName).
		Str("Chart", k.chartPath).
		Msg("Succesfully installed helm chart")
	return nil
}
