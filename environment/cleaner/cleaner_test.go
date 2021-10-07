package cleaner

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/smartcontractkit/integrations-framework/environment"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

func newNamespace(client *kubernetes.Clientset, labels map[string]string) (string, error) {
	ns, err := client.CoreV1().Namespaces().Create(
		context.Background(),
		&coreV1.Namespace{
			ObjectMeta: metaV1.ObjectMeta{
				GenerateName: "something-",
				Labels:       labels,
			},
		},
		metaV1.CreateOptions{},
	)
	if err != nil {
		return "", err
	}
	return ns.Name, err
}

func rmNamespace(client *kubernetes.Clientset, name string) error {
	if err := client.CoreV1().Namespaces().Delete(context.Background(), name, metaV1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func newNamespaces(number int, client *kubernetes.Clientset, labels map[string]string) error {
	for i := 0; i < number; i++ {
		if _, err := client.CoreV1().Namespaces().Create(
			context.Background(),
			&coreV1.Namespace{
				ObjectMeta: metaV1.ObjectMeta{
					GenerateName: "something-",
					Labels:       labels,
				},
			},
			metaV1.CreateOptions{},
		); err != nil {
			return err
		}
	}
	return nil
}

func getNamespace(client *kubernetes.Clientset, name string) (*coreV1.Namespace, error) {
	ns, err := client.CoreV1().Namespaces().Get(context.Background(), name, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return ns, nil
}

var _ = Describe("Environment Cleaner tests @cleaner", func() {
	var (
		clientset             *kubernetes.Clientset
		notATestNamespaceName string
		err                   error
	)
	Describe("Basic tests for timeout policy", func() {
		It("Must delete only labelled namespaces after timeout", func() {
			config, err := environment.K8sConfig()
			Expect(err).ShouldNot(HaveOccurred())
			clientset, err = kubernetes.NewForConfig(config)
			Expect(err).ShouldNot(HaveOccurred())
			c := NewCleaner(clientset, &Config{PollInterval: 1 * time.Second})
			// nolint
			go c.Run()
			notATestNamespaceName, err = newNamespace(clientset, map[string]string{
				"type": "not_a_test",
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = newNamespaces(3, clientset, map[string]string{
				"type":    "test",
				"policy":  "timeout",
				"timeout": "3s",
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = newNamespaces(2, clientset, map[string]string{
				"type":    "test",
				"policy":  "timeout",
				"timeout": "10s",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(func(g Gomega) int {
				ns, _ := clientset.CoreV1().Namespaces().List(context.Background(), metaV1.ListOptions{
					LabelSelector: environment.BasicTestNamespaceSelector,
				})
				return len(ns.Items)
			}, "20s", "1s").Should(Equal(2))
			Eventually(func(g Gomega) int {
				ns, _ := clientset.CoreV1().Namespaces().List(context.Background(), metaV1.ListOptions{
					LabelSelector: environment.BasicTestNamespaceSelector,
				})
				return len(ns.Items)
			}, "20s", "1s").Should(Equal(0))
			ns, _ := clientset.CoreV1().Namespaces().List(context.Background(), metaV1.ListOptions{
				LabelSelector: environment.BasicTestNamespaceSelector,
			})
			Expect(ns.Items).Should(HaveLen(0))
			n, err := getNamespace(clientset, notATestNamespaceName)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(n.Status.Phase)).Should(Equal("Active"))
		})
		AfterEach(func() {
			err = rmNamespace(clientset, notATestNamespaceName)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
