package Deployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/wait"
)

var _ = Describe("Check ClusterIssuers", func() {
	It("should find 'amq-ca-issuer' and 'amq-selfsigned-cluster-issuer' ClusterIssuers", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		expectedClusterIssuers := []string{"amq-ca-issuer", "amq-selfsigned-cluster-issuer"}

		// Use the REST client to get ClusterIssuers
		restClient := clientset.CoreV1().RESTClient()
		req := restClient.Get().
			Resource("clusterissuers").
			VersionedParams(nil, metav1.ParameterCodec)

		result := req.Do(context.TODO())
		Expect(result.Error()).To(BeNil(), "Error getting ClusterIssuers: %v", result.Error())

		clusterIssuersList, err := unstructuredListFromResponse(result.Raw())
		Expect(err).To(BeNil(), "Error converting result to UnstructuredList: %v", err)

		foundClusterIssuers := extractClusterIssuerNames(clusterIssuersList)

		// Check if all expected ClusterIssuers are found
		Expect(foundClusterIssuers).To(ConsistOf(expectedClusterIssuers),
			"Expected ClusterIssuers %v, but found %v", expectedClusterIssuers, foundClusterIssuers)
	})
})

func unstructuredListFromResponse(rawResponse []byte) (*unstructured.UnstructuredList, error) {
	ul := &unstructured.UnstructuredList{}
	err := ul.UnmarshalJSON(rawResponse)
	return ul, err
}

func extractClusterIssuerNames(clusterIssuersList *unstructured.UnstructuredList) []string {
	var foundClusterIssuers []string

	for _, item := range clusterIssuersList.Items {
		issuerName, found, _ := unstructured.NestedString(item.Object, "metadata", "name")
		if found {
			foundClusterIssuers = append(foundClusterIssuers, issuerName)
			fmt.Printf("Found ClusterIssuer: %s\n", issuerName)
		}
	}

	return foundClusterIssuers
}
