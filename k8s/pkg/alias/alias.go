package alias

import (
	"fmt"
	"strings"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/k8s/imports/k8s"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// ShortDur is a helper method for kube-janitor duration format
func ShortDur(d time.Duration) *string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return utils.Ptr(s)
}

func ConvertLabels(labels []string) (*map[string]*string, error) {
	cdk8sLabels := make(map[string]*string)
	for _, s := range labels {
		a := strings.Split(s, "=")
		if len(a) != 2 {
			return nil, fmt.Errorf("invalid label '%s' provided, please provide labels in format key=value", a)
		}
		cdk8sLabels[a[0]] = utils.Ptr(a[1])
	}
	return &cdk8sLabels, nil
}

// ConvertAnnotations converts a map[string]string to a *map[string]*string
func ConvertAnnotations(annotations map[string]string) *map[string]*string {
	a := make(map[string]*string)
	for k, v := range annotations {
		a[k] = utils.Ptr(v)
	}
	return &a
}

// EnvVarStr quick shortcut for string/string key/value var
func EnvVarStr(k, v string) *k8s.EnvVar {
	return &k8s.EnvVar{
		Name:  utils.Ptr(k),
		Value: utils.Ptr(v),
	}
}

// ContainerResources container resource requirements
func ContainerResources(reqCPU, reqMEM, limCPU, limMEM string) *k8s.ResourceRequirements {
	return &k8s.ResourceRequirements{
		Requests: &map[string]k8s.Quantity{
			"cpu":    k8s.Quantity_FromString(utils.Ptr(reqCPU)),
			"memory": k8s.Quantity_FromString(utils.Ptr(reqMEM)),
		},
		Limits: &map[string]k8s.Quantity{
			"cpu":    k8s.Quantity_FromString(utils.Ptr(limCPU)),
			"memory": k8s.Quantity_FromString(utils.Ptr(limMEM)),
		},
	}
}
