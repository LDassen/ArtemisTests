package Deployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	certmanagerv1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	"github.com/jetstack/cert-manager/pkg/client/clientset/versioned"
)

var _ = Describe("Check if Certificates are present in the correct namespace", func() {
	It("should ensure the certificates exist in the specified namespaces", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		certManagerClient, err := versioned.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Cert-Manager client: %v", err)

		// Check for the certificate in 'activemq-artemis-brokers' namespace
		amqTLSCertName := "amq-tls-acceptor-cert"
		amqTLSCertNamespace := "activemq-artemis-brokers"
		_, err = certManagerClient.CertmanagerV1().Certificates(amqTLSCertNamespace).Get(context.TODO(), amqTLSCertName, metav1.GetOptions{})
		Expect(err).To(BeNil(), "Error getting certificate '%s' in namespace '%s': %v", amqTLSCertName, amqTLSCertNamespace, err)
		fmt.Printf("Certificate '%s' found in namespace '%s'\n", amqTLSCertName, amqTLSCertNamespace)

		// Check for the certificate in 'cert-manager' namespace
		selfSignedCACertName := "amq-selfsigned-ca"
		selfSignedCACertNamespace := "cert-manager"
		_, err = certManagerClient.CertmanagerV1().Certificates(selfSignedCACertNamespace).Get(context.TODO(), selfSignedCACertName, metav1.GetOptions{})
		Expect(err).To(BeNil(), "Error getting certificate '%s' in namespace '%s': %v", selfSignedCACertName, selfSignedCACertNamespace, err)
		fmt.Printf("Certificate '%s' found in namespace '%s'\n", selfSignedCACertName, selfSignedCACertNamespace)
	})
})
