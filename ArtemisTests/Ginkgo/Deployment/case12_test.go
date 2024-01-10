package Deployment_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		// Check 'amq-ca-issuer'
		issuer1, err := clientset.CertmanagerV1().ClusterIssuers().Get(context.TODO(), "amq-ca-issuer", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred(), "Error while checking 'amq-ca-issuer'")
		Expect(issuer1.Spec.Acme).To(BeNil()) // Assuming you are not using ACME, adjust as needed
		Expect(issuer1.Status.Conditions[0].Type).To(Equal("Ready"))
		Expect(issuer1.Status.Conditions[0].Status).To(Equal("True"))

		// Check 'amq-selfsigned-cluster-issuer'
		issuer2, err := clientset.CertmanagerV1().ClusterIssuers().Get(context.TODO(), "amq-selfsigned-cluster-issuer", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred(), "Error while checking 'amq-selfsigned-cluster-issuer'")
		Expect(issuer2.Spec.Acme).To(BeNil()) // Assuming you are not using ACME, adjust as needed
		Expect(issuer2.Status.Conditions[0].Type).To(Equal("Ready"))
		Expect(issuer2.Status.Conditions[0].Status).To(Equal("True"))
	})
})
