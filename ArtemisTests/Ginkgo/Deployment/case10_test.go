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

var _ = Describe("Check ClusterIssuers", func() {
	It("should find 'amq-ca-issuer' and 'amq-selfsigned-cluster-issuer' ClusterIssuers with ready set to true in the 'cert-manager' namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "cert-manager"
		expectedClusterIssuers := []string{"amq-ca-issuer", "amq-selfsigned-cluster-issuer"}

		clusterIssuers, err := clientset.
			CertificatesV1().
			ClusterIssuers().
			List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error getting ClusterIssuers: %v", err)

		var foundClusterIssuers []string
		for _, issuer := range clusterIssuers.Items {
			if issuer.Status.Conditions.IsTrueFor("Ready") {
				foundClusterIssuers = append(foundClusterIssuers, issuer.Name)
				fmt.Printf("Found ready ClusterIssuer: %s\n", issuer.Name)
			}
		}

		// Check if all expected ClusterIssuers are found
		Expect(foundClusterIssuers).To(ConsistOf(expectedClusterIssuers),
			"Expected ready ClusterIssuers %v, but found %v", expectedClusterIssuers, foundClusterIssuers)
	})
})
