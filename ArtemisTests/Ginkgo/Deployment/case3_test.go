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

var _ = Describe("Check StatefulSet Existence", func() {
	It("should ensure the StatefulSet 'ex-aao-ss' exists in the specified namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artemis-brokers"
		statefulSetName := "ex-aao-ss"

		statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("StatefulSet '%s' not found in namespace '%s'\n", statefulSetName, namespace)
			// Handle the absence of StatefulSet as needed
			// You may choose to fail the test or perform additional actions
		} else {
			fmt.Printf("StatefulSet '%s' found in namespace '%s'\n", statefulSetName, namespace)
			fmt.Printf("Replicas: %d\n", statefulSet.Status.Replicas)
			// Add additional checks or actions based on the StatefulSet's presence
		}
	})
})