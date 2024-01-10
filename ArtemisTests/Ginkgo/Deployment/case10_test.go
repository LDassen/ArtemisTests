package Deployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/util/jsonpath"
	"k8s.io/client-go/util/wait"
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

		var clusterIssuersList runtime.Object
		err = json.NewYAMLSerializer(json.DefaultMetaFactory, nil, nil).
			Decode(result.Body(), &clusterIssuersList)
		Expect(err).To(BeNil(), "Error converting result to List: %v", err)

		foundClusterIssuers := extractClusterIssuerNames(clusterIssuersList)

		// Check if all expected ClusterIssuers are found
		Expect(foundClusterIssuers).To(ConsistOf(expectedClusterIssuers),
			"Expected ready ClusterIssuers %v, but found %v", expectedClusterIssuers, foundClusterIssuers)
	})
})

func extractClusterIssuerNames(clusterIssuersList runtime.Object) []string {
	var foundClusterIssuers []string

	jsonPath := jsonpath.New("clusterissuer-name")
	err := jsonPath.Parse("{.items[*].metadata.name}")
	if err != nil {
		Fail(fmt.Sprintf("Error parsing JSON path: %v", err))
	}

	err = jsonPath.Execute(wait.NeverStop, result.Body(), &foundClusterIssuers)
	if err != nil {
		Fail(fmt.Sprintf("Error extracting ClusterIssuer names: %v", err))
	}

	return foundClusterIssuers
}
