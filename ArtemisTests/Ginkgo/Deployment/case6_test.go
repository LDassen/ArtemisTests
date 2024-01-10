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

var _ = Describe("Check ConfigMap Existence", func() {
	It("should ensure the ConfigMap 'ca-bundle' exists in the specified namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artemis-brokers"
		configMapName := "ca-bundle"

		_, err = clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configMapName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("ConfigMap '%s' not found in namespace '%s'\n", configMapName, namespace)
			Expect(err).ToNot(HaveOccurred(), "Expected ConfigMap to be missing, but got an error.")
		} else {
			fmt.Printf("ConfigMap '%s' found in namespace '%s'\n", configMapName, namespace)
			Expect(err).To(BeNil(), "Expected ConfigMap to exist, but got an error.")
		}
	})
})
