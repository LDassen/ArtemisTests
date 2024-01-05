package SSLConfig_test

import (
	"context"

	certmanager "github.com/cert-manager/cert-manager"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	cmmetav1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ = ginkgo.Describe("Kafka Certificates and Secrets", func() {
	var clientset *kubernetes.Clientset
	var certManagerClientset *certmanager.Clientset

	ginkgo.BeforeEach(func() {
		// Set up the Kubernetes client
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		clientset, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Set up the cert-manager client
		certManagerClientset, err = certmanager.NewForConfig(config)
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
				cert, err := certManagerClientset.CertmanagerV1().Certificates(namespace).Get(context.TODO(), certName, metav1.GetOptions{})
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(cert.Status.Conditions).To(gomega.ContainElement(cmmetav1.Condition{
					Type:   cmmetav1.CertificateConditionReady,
					Status: cmmetav1.ConditionTrue,
				}))
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
