package Deployment_test

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
	It("should check cluster secrets", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		// case 7
		checkSecretInNamespace(clientset, "activemq-artemis-brokers", "amq-ssl-secret", 1, "Check amq-ssl-secret in activemq-artemis-brokers namespace")
		checkSecretInNamespace(clientset, "activemq-artemis-brokers", "ex-aao-credentials-secret", 1, "Check ex-aao-credentials-secret in activemq-artemis-brokers namespace")
		checkSecretInNamespace(clientset, "activemq-artemis-brokers", "ex-aao-netty-secret", 1, "Check ex-aao-netty-secret in activemq-artemis-brokers namespace")
		checkSecretInNamespace(clientset, "activemq-artemis-brokers", "ex-aao-props", 1, "Check ex-aao-props in activemq-artemis-brokers namespace")
		checkSecretInNamespace(clientset, "activemq-artemis-brokers", "jks-password-secret", 1, "Check jks-password-secret in activemq-artemis-brokers namespace")
		checkSecretInNamespace(clientset, "activemq-artemis-brokers", "secret-security-ex-aao-prop", 1, "Check secret-security-ex-aao-prop in activemq-artemis-brokers namespace")
		checkSecretInNamespace(clientset, "activemq-artemis-brokers", "ssl-acceptor-ssl-secret", 1, "Check ssl-acceptor-ssl-secret in activemq-artemis-brokers namespace")

		// case 8
		checkSecretInNamespace(clientset, "activemq-artemis-operator", "amq-ssl-secret", 1, "Check amq-ssl-secret in activemq-artemis-operator namespace")
		checkSecretInNamespace(clientset, "activemq-artemis-operator", "jks-password-secret", 1, "Check jks-password-secret in activemq-artemis-operator namespace")

		// case 9
		checkSecretInNamespace(clientset, "cert-manager", "root-secret", 1, "Check root-secret in cert-manager namespace")
	})
})

func checkSecretInNamespace(clientset *kubernetes.Clientset, namespace, secretName string, expectedCount int, description string) {
	secrets, err := clientset.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", secretName),
	})
	Expect(err).To(BeNil(), "Error getting secrets in namespace %s: %v", namespace, err)

	actualCount := len(secrets.Items)
	Expect(actualCount).To(Equal(expectedCount), "%s - Expected %d '%s' secret(s) in namespace %s, but found %d", description, expectedCount, secretName, namespace, actualCount)

	if actualCount > 0 {
		fmt.Printf("Found '%s' secret(s) in namespace %s\n", secretName, namespace)
	}
}
