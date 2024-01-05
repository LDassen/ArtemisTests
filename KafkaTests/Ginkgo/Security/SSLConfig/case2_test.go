package SSLConfig_test

import (
	"context"
	"time"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ = ginkgo.Describe("Kafka Certificates and Secrets", func() {
	var clientset *kubernetes.Clientset

	ginkgo.BeforeEach(func() {
		// Set up the Kubernetes client
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		clientset, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.Context("in kafka-brokers namespace", func() {
		namespace := "kafka-brokers"

		ginkgo.It("should have two ready certificates", func() {
			certNames := []string{
				"kafka-brokers-controller.kafka-brokers.mgt.cluster.local",
				"kafka-brokers-headless.kafka-brokers.svc.cluster.local",
			}

			for _, certName := range certNames {
				// Assuming the use of cert-manager, change this according to your certificate API
				cert, err := clientset.CertificatesV1().Certificates(namespace).Get(context.TODO(), certName, metav1.GetOptions{})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(cert.Status.Conditions).To(gomega.ContainElement( /* check for Ready status condition */ ))
			}
		})

		ginkgo.It("should have specific secrets", func() {
			secretNames := []string{
				"kafka-brokers-controller",
				"kafka-brokers-server-certificate",
			}

			for _, secretName := range secretNames {
				_, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			}
		})
	})
})
