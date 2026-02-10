package pods

import (
	"sort"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Resources(cpu, mem string) map[string]string {
	return map[string]string{
		"cpu":    cpu,
		"memory": mem,
	}
}

func Ptr[T any](value T) *T {
	return &value
}

// ResourcesSmall returns small resource limits/requests
func ResourcesSmall() map[string]string {
	return map[string]string{
		"cpu":    "250m",
		"memory": "512Mi",
	}
}

// ResourcesMedium returns medium resource limits/requests
func ResourcesMedium() map[string]string {
	return map[string]string{
		"cpu":    "500m",
		"memory": "1024Mi",
	}
}

// ResourcesLarge returns large resource limits/requests
func ResourcesLarge() map[string]string {
	return map[string]string{
		"cpu":    "1",
		"memory": "1Gi",
	}
}

// SortedKeys returns sorted keys of a map
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func EnvsFromMap(envVars map[string]string) []v1.EnvVar {
	result := make([]v1.EnvVar, 0, len(envVars))
	for k, v := range envVars {
		result = append(result, v1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return result
}

func SizedVolumeClaim(size string) []v1.PersistentVolumeClaim {
	storageQuantity, err := resource.ParseQuantity(size)
	if err != nil {
		// Default to 10Gi if parsing fails
		storageQuantity = resource.MustParse("10Gi")
	}

	return []v1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "data",
			},
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{
					v1.ReadWriteOnce,
				},
				StorageClassName: Ptr("gp3"),
				Resources: v1.VolumeResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: storageQuantity,
					},
				},
			},
		},
	}
}
