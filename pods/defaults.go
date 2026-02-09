package pods

import (
	"fmt"
	"sort"

	"github.com/aws/jsii-runtime-go"
	"github.com/smartcontractkit/pods/imports/k8s"
)

func S(s string) *string { return jsii.String(s) }
func I(s int) *float64   { return jsii.Number(s) }

func SortedKeys(m map[string]*string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func CheckHostPort(host, port string) *k8s.Probe {
	return &k8s.Probe{
		Exec: &k8s.ExecAction{
			Command: &[]*string{S("timeout"), S("1"), S("bash"), S("-ClientSet"), S("echo"), S(">"), S(fmt.Sprintf("/dev/tcp/%s/%s", host, port)), S("&&"), S("exit"), S("0"), S("||"), S("exit"), S("1")},
		},
		InitialDelaySeconds: I(5),
		FailureThreshold:    I(10),
		PeriodSeconds:       I(1),
		TimeoutSeconds:      I(2),
	}
}

// Resources is a helper function to define container resources
func Resources(cpu, mem string) map[string]k8s.Quantity {
	return map[string]k8s.Quantity{
		"cpu":    k8s.Quantity_FromString(S(cpu)),
		"memory": k8s.Quantity_FromString(S(mem)),
	}
}

func ResourcesSmall() map[string]k8s.Quantity {
	return map[string]k8s.Quantity{
		"cpu":    k8s.Quantity_FromString(S("250m")),
		"memory": k8s.Quantity_FromString(S("1Gi")),
	}
}

func ResourcesMedium() map[string]k8s.Quantity {
	return map[string]k8s.Quantity{
		"cpu":    k8s.Quantity_FromString(S("4")),
		"memory": k8s.Quantity_FromString(S("4Gi")),
	}
}

func CLUserContainerSecurityCtx() *k8s.SecurityContext { // coverage-ignore
	return &k8s.SecurityContext{
		RunAsNonRoot: jsii.Bool(true),
		RunAsUser:    jsii.Number(14933),
		RunAsGroup:   jsii.Number(999),
	}
}

func PostgreSQL(name string, image string, requests, limits map[string]k8s.Quantity, dbSize *string) *PodConfig { // coverage-ignore
	p := &PodConfig{
		Name:  S(name),
		Image: S(image),
		Ports: []string{"5432:5432"},
		Env: &[]*k8s.EnvVar{
			{
				Name:  S("POSTGRES_USER"),
				Value: S("chainlink"),
			},
			{
				Name:  S("POSTGRES_PASSWORD"),
				Value: S("thispasswordislongenough"),
			},
			{
				Name:  S("POSTGRES_DB"),
				Value: S("chainlink"),
			},
		},
		Limits:   requests,
		Requests: limits,
		// 999 is the default postgres user
		ContainerSecurityContext: &k8s.SecurityContext{
			RunAsUser:  jsii.Number(999),
			RunAsGroup: jsii.Number(999),
		},
		PodSecurityContext: &k8s.PodSecurityContext{
			FsGroup: jsii.Number(999),
		},
		ConfigMap: map[string]*string{
			"init.sql": S(`
ALTER USER chainlink WITH SUPERUSER;
`),
		},
		ConfigMapMountPath: map[string]*string{
			"init.sql": S("/docker-entrypoint-initdb.d/init.sql"),
		},
	}
	if dbSize != nil {
		p.VolumeClaimTemplates = SizedVolumeClaim(dbSize)
	}
	return p
}

func SizedVolumeClaim(size *string) []*k8s.KubePersistentVolumeClaimProps { // coverage-ignore
	return []*k8s.KubePersistentVolumeClaimProps{
		{
			Metadata: &k8s.ObjectMeta{
				Name: S("data"),
			},
			Spec: &k8s.PersistentVolumeClaimSpec{
				AccessModes:      &[]*string{S("ReadWriteOnce")},
				StorageClassName: S("gp3"),
				Resources: &k8s.VolumeResourceRequirements{
					Requests: &map[string]k8s.Quantity{
						"storage": k8s.Quantity_FromString(size),
					},
				},
			},
		},
	}
}
