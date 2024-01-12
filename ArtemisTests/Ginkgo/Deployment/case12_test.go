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

var _ = Describe("Check ClusterIssuer Existence", func() {
	It("should ensure the ClusterIssuer 'amq-ca-issuer' exists in the 'activemq-artemis-brokers' namespace", func() {
		// Get in-cluster config
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		// Create Kubernetes client
		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		// Specify namespace and ClusterIssuer name
		namespace := "activemq-artemis-brokers"
		issuerName := "amq-ca-issuer"

		// Check if the ClusterIssuer exists
		_, err = clientset.AdmissionregistrationV1().ClusterIssuers().Get(context.TODO(), issuerName, metav1.GetOptions{})
		if err != nil {
			// ClusterIssuer not found
			fmt.Printf("ClusterIssuer '%s' not found in namespace '%s'\n", issuerName, namespace)
			Expect(err).ToNot(HaveOccurred(), "Expected ClusterIssuer to be missing, but got an error.")
		} else {
			// ClusterIssuer found
			fmt.Printf("ClusterIssuer '%s' found in namespace '%s'\n", issuerName, namespace)
			Expect(err).To(BeNil(), "Expected ClusterIssuer to exist, but got an error.")
		}
	})

	It("should ensure the ClusterIssuer 'root-secret' exists in the 'cert-manager' namespace", func() {
		// Get in-cluster config
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		// Create Kubernetes client
		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		// Specify namespace and ClusterIssuer name
		namespace := "cert-manager"
		issuerName := "root-secret"

		// Check if the ClusterIssuer exists
		_, err = clientset.AdmissionregistrationV1().ClusterIssuers().Get(context.TODO(), issuerName, metav1.GetOptions{})
		if err != nil {
			// ClusterIssuer not found
			fmt.Printf("ClusterIssuer '%s' not found in namespace '%s'\n", issuerName, namespace)
			Expect(err).ToNot(HaveOccurred(), "Expected ClusterIssuer to be missing, but got an error.")
		} else {
			// ClusterIssuer found
			fmt.Printf("ClusterIssuer '%s' found in namespace '%s'\n", issuerName, namespace)
			Expect(err).To(BeNil(), "Expected ClusterIssuer to exist, but got an error.")
		}
	})
})
