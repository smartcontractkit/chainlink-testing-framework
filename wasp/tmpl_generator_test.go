package wasp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// use
// kubectl -n default get po --output=json | jq .items
// to write tests for real-world cases

func TestGenerateDataPositive(t *testing.T) {
	tests := []struct {
		name            string
		namespace       string
		mockPods        *corev1.PodList
		mockError       error
		includeWorkload bool
		expectedCases   int
	}{
		{
			name:      "single pod with correct and unique instance label",
			namespace: "test-namespace",
			mockPods: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-pod-1",
							Namespace: "test-namespace",
							Labels: map[string]string{
								DefaultUniqLabel: "app-1",
							},
						},
					},
				},
			},
			includeWorkload: false,
			expectedCases:   2,
		},
		{
			name:      "single pod with correct and unique instance label + workload",
			namespace: "test-namespace",
			mockPods: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-pod-1",
							Namespace: "test-namespace",
							Labels: map[string]string{
								DefaultUniqLabel: "app-1",
							},
						},
					},
				},
			},
			includeWorkload: true,
			expectedCases:   2,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockK8s{
				Pods: tt.mockPods,
				Err:  tt.mockError,
			}
			generator, err := NewLoadTestGenBuilder(mockClient, tt.namespace).
				Workload(tt.includeWorkload).Build()
			require.NoError(t, err)

			err = generator.Read()
			require.NoError(t, err)
			
			testCases, err := generator.GenerateTestCases()
			require.NoError(t, err)
			assert.Len(t, testCases, tt.expectedCases)

			experiments, err := generator.GenerateTableTest()
			require.NoError(t, err)
			assert.Contains(t, experiments, "package main")
			assert.Contains(t, experiments, tt.namespace)
			assert.Contains(t, experiments, "TestGeneratedLoadChaos")

			if tt.includeWorkload {
				assert.Contains(t, experiments, "wasp.NewGenerator")
				assert.Contains(t, experiments, "ExampleHTTPGun")
			} else {
				assert.NotContains(t, experiments, "wasp.NewGenerator")
			}
		})
	}
}

func TestGenerateDataNegative(t *testing.T) {
	tests := []struct {
		name            string
		namespace       string
		mockPods        *corev1.PodList
		mockError       error
		includeWorkload bool
		expectedCases   int
		expectError     bool
	}{
		{
			name:      "single pod but incorrect instance annotation",
			namespace: "test-namespace",
			mockPods: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-pod-1",
							Namespace: "test-namespace",
							Labels: map[string]string{
								"app": "app1",
							},
						},
					},
				},
			},
			includeWorkload: false,
			expectedCases:   0,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockK8s{
				Pods: tt.mockPods,
				Err:  tt.mockError,
			}
			generator, err := NewLoadTestGenBuilder(mockClient, tt.namespace).
				Workload(tt.includeWorkload).Build()
			require.NoError(t, err)

			err = generator.Read()
			require.NoError(t, err)

			experiments, err := generator.GenerateTableTest()
			require.NoError(t, err)
			assert.Contains(t, experiments, "package main")
			assert.Contains(t, experiments, tt.namespace)
			assert.Contains(t, experiments, "TestGeneratedLoadChaos")

			if tt.includeWorkload {
				assert.Contains(t, experiments, "wasp.NewGenerator")
				assert.Contains(t, experiments, "ExampleGun")
			} else {
				assert.NotContains(t, experiments, "wasp.NewGenerator")
			}
		})
	}
}

// func TestGenerateFiles(t *testing.T) {
// 	t.Skip("it's fine to fail until the lib is merged because 'replace' directive is needed but can't be merged")
// 	tests := []struct {
// 		name            string
// 		includeWorkload bool
// 		expectGunFile   bool
// 	}{
// 		{
// 			name:            "with workload generation",
// 			includeWorkload: true,
// 			expectGunFile:   true,
// 		},
// 		{
// 			name:            "without workload generation",
// 			includeWorkload: false,
// 			expectGunFile:   false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockClient := &MockK8s{
// 				Pods: &corev1.PodList{
// 					Items: []corev1.Pod{
// 						{
// 							ObjectMeta: metav1.ObjectMeta{
// 								Name:      "test-pod",
// 								Namespace: "test-namespace",
// 								Labels: map[string]string{
// 									"app": "test-app",
// 								},
// 							},
// 						},
// 					},
// 				},
// 			}
// 			tmpDir, err := os.MkdirTemp("", "chaos-test-*")
// 			require.NoError(t, err)
// 			defer os.RemoveAll(tmpDir)

// 			generator, err := NewLoadTestGenBuilder(mockClient, "test-namespace").
// 				Workload(tt.includeWorkload).
// 				OutputDir(tmpDir).
// 				GoModName("github.com/test/chaos-tests").
// 				Build()
// 			require.NoError(t, err)

// 			require.NoError(t, generator.Read())
// 			require.NoError(t, generator.Write())

// 			expectedFiles := []string{"go.mod", "chaos_test.go"}

// 			for _, file := range expectedFiles {
// 				filePath := filepath.Join(tmpDir, file)
// 				_, err := os.Stat(filePath)
// 				assert.NoError(t, err, "file %s should exist", file)
// 			}

// 			// Verify file contents, can be tidied and builded
// 			goModPath := filepath.Join(tmpDir, "go.mod")
// 			goModContent, err := os.ReadFile(goModPath)
// 			require.NoError(t, err)
// 			assert.Contains(t, string(goModContent), "module github.com/test/chaos-tests")
// 			assert.Contains(t, string(goModContent), "go 1.25")

// 			testPath := filepath.Join(tmpDir, "chaos_test.go")
// 			testContent, err := os.ReadFile(testPath)
// 			require.NoError(t, err)

// 			contentStr := string(testContent)
// 			assert.Contains(t, contentStr, "package main")
// 			assert.Contains(t, contentStr, "TestGeneratedLoadChaos")
// 			assert.Contains(t, contentStr, "Fail pod test-app")
// 			assert.Contains(t, contentStr, "Network delay for test-app")

// 			if tt.includeWorkload {
// 				assert.Contains(t, contentStr, "ExampleGun")
// 				assert.Contains(t, contentStr, "wasp.NewGenerator")
// 			}
// 		})
// 	}
// }
