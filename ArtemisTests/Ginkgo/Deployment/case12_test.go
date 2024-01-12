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

var _ = Describe("Check ClusterIssuers Existence", func() {
	It("should ensure ClusterIssuers exist in the specified namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "cert-manager"

		// List all ClusterIssuers in the namespace
		clusterIssuersList, err := clientset.AdmissionregistrationV1().ClusterIssuers().List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error listing ClusterIssuers: %v", err)

		// Names of ClusterIssuers to find
		clusterIssuerNames := []string{"amq-ca-issuer", "amq-selfsigned-cluster-issuer"}

		// Check each ClusterIssuer's existence
		for _, clusterIssuerName := range clusterIssuerNames {
			found := false
			for _, ci := range clusterIssuersList.Items {
				if ci.Name == clusterIssuerName {
					found = true
					fmt.Printf("ClusterIssuer '%s' found in namespace '%s'\n", ci.Name, namespace)
					break
				}
			}
			Expect(found).To(BeTrue(), "ClusterIssuer '%s' not found in namespace '%s'", clusterIssuerName, namespace)
		}
	})
})
