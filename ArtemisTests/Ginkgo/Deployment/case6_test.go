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

var _ = Describe("Check if ca-bundle ConfigMap is synced", func() {
	It("should ensure ca-bundle ConfigMap is synced in the cert-manager namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "cert-manager"
		configMapName := "ca-bundle"

		configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
		Expect(err).To(BeNil(), "Error getting ConfigMap '%s' in namespace '%s': %v", configMapName, namespace, err)

		// Check if 'synced' field exists and is set to 'true'
		synced, found := configMap.Data["synced"]
		Expect(found).To(BeTrue(), "Field 'synced' not found in ConfigMap '%s' in namespace '%s'", configMapName, namespace)
		Expect(synced).To(Equal("true"), "Expected 'synced' to be 'true' in ConfigMap '%s' in namespace '%s', but found '%s'", configMapName, namespace, synced)
	})
})