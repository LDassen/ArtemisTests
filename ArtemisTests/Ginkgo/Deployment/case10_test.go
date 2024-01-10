package Deployment_test

import (
	"context"
	"fmt"
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Check ClusterIssuers", func() {
	It("should find 'amq-ca-issuer' and 'amq-selfsigned-cluster-issuer' ClusterIssuers with ready set to true", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		expectedClusterIssuers := []string{"amq-ca-issuer", "amq-selfsigned-cluster-issuer"}

		// Use the REST client to get ClusterIssuers
		restClient := clientset.CoreV1().RESTClient()
		req := restClient.Get().
			Resource("clusterissuers").
			VersionedParams(&metav1.ListOptions{TypeMeta: metav1.TypeMeta{Kind: "ClusterIssuer"}}, metav1.ParameterCodec)

		result := req.Do(context.TODO())
		Expect(result.Error()).To(BeNil(), "Error getting ClusterIssuers: %v", result.Error())

		var clusterIssuersList metav1.List
		err = result.Into(&clusterIssuersList)
		Expect(err).To(BeNil(), "Error converting result to List: %v", err)

		var foundClusterIssuers []string
		for _, item := range clusterIssuersList.Items {
			issuer, ok := item.(*unstructured.Unstructured)
			if !ok {
				Fail("Failed to convert item to Unstructured")
			}

			issuerName, found, _ := unstructured.NestedString(issuer.Object, "metadata", "name")
			if !found {
				Fail("Failed to extract issuer name")
			}

			conditions, found, _ := unstructured.NestedSlice(issuer.Object, "status", "conditions")
			if !found {
				Fail("Failed to extract conditions")
			}

			// Check if 'Ready' condition is true
			ready := false
			for _, condition := range conditions {
				status, found, _ := unstructured.NestedString(condition.(map[string]interface{}), "status")
				if found && status == "True" {
					ready = true
					break
				}
			}

			if ready {
				foundClusterIssuers = append(foundClusterIssuers, issuerName)
				fmt.Printf("Found ready ClusterIssuer: %s\n", issuerName)
			}
		}

		// Check if all expected ClusterIssuers are found
		Expect(foundClusterIssuers).To(ConsistOf(expectedClusterIssuers),
			"Expected ready ClusterIssuers %v, but found %v", expectedClusterIssuers, foundClusterIssuers)
	})
})
