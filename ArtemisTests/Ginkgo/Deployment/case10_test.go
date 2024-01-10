package Deployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/jsonpath"
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
		req := restClient.
			Get().
			Resource("clusterissuers").
			VersionedParams(&metav1.ListOptions{TypeMeta: metav1.TypeMeta{Kind: "ClusterIssuer"}}, metav1.ParameterCodec)

		result := req.Do(context.TODO())
		Expect(result.Error()).To(BeNil(), "Error getting ClusterIssuers: %v", result.Error())

		// Check if the response status code indicates success
		Expect(result.StatusCode()).To(Equal(200), "Unexpected status code: %d", result.StatusCode())

		// Parse the JSONPath template
		template := "{range .items[*]}{.metadata.name}{end}"
		parser := jsonpath.New("clusterissuer-name")
		parser.AllowMissingKeys(true)
		err = parser.Parse(template)
		Expect(err).To(BeNil(), "Error parsing JSONPath template: %v", err)

		// Evaluate JSONPath template on the response object
		var foundClusterIssuers []string
		err = parser.Execute(result, &foundClusterIssuers)
		Expect(err).To(BeNil(), "Error evaluating JSONPath: %v", err)

		// Check if all expected ClusterIssuers are found
		Expect(foundClusterIssuers).To(ConsistOf(expectedClusterIssuers),
			"Expected ready ClusterIssuers %v, but found %v", expectedClusterIssuers, foundClusterIssuers)
	})
})
