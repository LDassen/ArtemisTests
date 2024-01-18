package CheckDeployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Secrets Check", func() {
	It("should check namespace secrets", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		// Checking secrets in the dev-vault-1 namespace
		checkSecret(clientset, "dev-vault-1", "dev-vault-1-configurer")
		checkSecret(clientset, "dev-vault-1", "dev-vault-1-raw-config")
		checkSecret(clientset, "dev-vault-1", "dev-vault-1-tls")
		checkSecret(clientset, "dev-vault-1", "dev-vault-credentials-1")
		checkSecret(clientset, "dev-vault-1", "vault-unseal-keys")
	})
})

func checkSecret(clientset *kubernetes.Clientset, namespace, secretName string) {
	_, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	Expect(err).To(BeNil(), "Error getting secret '%s' in namespace '%s': %v", secretName, namespace, err)
	fmt.Printf("Secret '%s' found in namespace '%s'\n", secretName, namespace)
}

