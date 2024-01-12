package Deployment_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	certmanagerv1 "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned/"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
)

var _ = Describe("ClusterIssuers Check", func() {
	var clientset *kubernetes.Clientset
	var certManagerClientset *certmanagerv1.CertmanagerV1Client

	BeforeEach(func() {
		// Set up the Kubernetes client
		config, err := rest.InClusterConfig()
		Expect(err).NotTo(HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		// Set up the Cert-Manager client
		certManagerClientset, err = certmanagerv1.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should find 'amq-ca-issuer' and 'amq-selfsigned-cluster-issuer'", func() {
		// Check 'amq-ca-issuer'
		issuer1, err := certManagerClientset.ClusterIssuers("cert-manager.io").Get(context.TODO(), "amq-ca-issuer", metav1.GetOptions{})
		if err == nil {
			println("amq-ca-issuer exists:")
			println("Ready:", issuer1.Status.Conditions[0].Status)
		} else {
			println("amq-ca-issuer does not exist")
		}

		// Check 'amq-selfsigned-cluster-issuer'
		issuer2, err := certManagerClientset.ClusterIssuers("cert-manager.io").Get(context.TODO(), "amq-selfsigned-cluster-issuer", metav1.GetOptions{})
		if err == nil {
			println("amq-selfsigned-cluster-issuer exists:")
			println("Ready:", issuer2.Status.Conditions[0].Status)
		} else {
			println("amq-selfsigned-cluster-issuer does not exist")
		}
	})
})