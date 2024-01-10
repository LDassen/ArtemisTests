package Deployment_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	certv1 "k8s.io/api/certificates/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
)

var _ = Describe("ClusterIssuers Check", func() {
	var clientset *kubernetes.Clientset

	BeforeEach(func() {
		// Set up the Kubernetes client
		config, err := rest.InClusterConfig()
		Expect(err).NotTo(HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should find 'amq-ca-issuer' and 'amq-selfsigned-cluster-issuer'", func() {
		// Check for the presence of 'amq-ca-issuer'
		_, err := clientset.CertificatesV1().CertificateSigningRequests().Get(context.TODO(), "amq-ca-issuer", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred(), "Error while checking 'amq-ca-issuer'")

		// Check for the presence of 'amq-selfsigned-cluster-issuer'
		_, err = clientset.CertificatesV1().CertificateSigningRequests().Get(context.TODO(), "amq-selfsigned-cluster-issuer", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred(), "Error while checking 'amq-selfsigned-cluster-issuer'")
	})
})
