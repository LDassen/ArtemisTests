package Deployment_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func TestDeployment(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Deployment Suite")
}

var _ = Describe("Check ClusterIssuers Existence", func() {
	It("should ensure ClusterIssuers exist in the specified namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		dynamicClient, err := dynamic.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating dynamic Kubernetes client: %v", err)

		namespace := "cert-manager"
		resourceType := "clusterissuers"

		// List all ClusterIssuers in the namespace
		resourceList, err := dynamicClient.Resource(
			&unstructured.Unstructured{Object: map[string]interface{}{"apiVersion": "cert-manager.io/v1", "kind": "ClusterIssuer"}},
		).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error listing ClusterIssuers: %v", err)

		// Names of ClusterIssuers to find
		clusterIssuerNames := []string{"amq-ca-issuer", "amq-selfsigned-cluster-issuer"}

		// Check each ClusterIssuer's existence and readiness
		for _, clusterIssuerName := range clusterIssuerNames {
			found := false
			for _, item := range resourceList.Items {
				if item.GetName() == clusterIssuerName {
					found = true
					fmt.Printf("ClusterIssuer '%s' found in namespace '%s'\n", clusterIssuerName, namespace)

					// Perform additional checks if needed
					// Check the conditions
					conditions, _, _ := unstructured.NestedSlice(item.Object, "status", "conditions")
					Expect(conditions).To(HaveLen(1), "Expected ClusterIssuer to have one condition.")
					Expect(conditions[0].(map[string]interface{})["type"]).To(Equal("Ready"), "Expected ClusterIssuer condition to be Ready.")
					Expect(conditions[0].(map[string]interface{})["status"]).To(Equal("True"), "Expected ClusterIssuer condition status to be True.")
					break
				}
			}
			Expect(found).To(BeTrue(), "ClusterIssuer '%s' not found in namespace '%s'", clusterIssuerName, namespace)
		}
	})
})
