package Deployment_test

import (
	"context"
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Check if ca-bundle ConfigMap is synced", func() {
	It("should ensure ca-bundle ConfigMap is synced", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		configMapName := "ca-bundle"

		// Specify an empty string "" to represent the cluster-wide search
		configMap, err := clientset.CoreV1().ConfigMaps("").Get(context.TODO(), configMapName, metav1.GetOptions{})
		Expect(err).To(BeNil(), "Error getting ConfigMap '%s': %v", configMapName, err)

		// Assuming 'Bundles' is a JSON-encoded string
		var bundlesData map[string]interface{}
		Expect(json.Unmarshal([]byte(configMap.Data["Bundles"]), &bundlesData)).To(BeNil(), "Error decoding 'Bundles' field in ConfigMap '%s'", configMapName)

		// Check if 'status.conditions' array exists
		conditions, found := bundlesData["status"].(map[string]interface{})["conditions"].([]interface{})
		Expect(found).To(BeTrue(), "Field 'status.conditions' not found in 'Bundles' of ConfigMap '%s'", configMapName)

		// Check if any condition has 'type' set to 'Synced' and 'status' set to 'True'
		var syncedConditionFound bool
		for _, condition := range conditions {
			conditionMap := condition.(map[string]interface{})
			if conditionMap["type"] == "Synced" && conditionMap["status"] == "True" {
				syncedConditionFound = true
				break
			}
		}

		Expect(syncedConditionFound).To(BeTrue(), "Expected 'Synced' condition with 'status' 'True' in ConfigMap '%s'", configMapName)
	})
})
