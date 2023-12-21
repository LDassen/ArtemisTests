package Checkpods_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Artemis Broker Pods", func() {
	It("should have the correct number of 'broker' pods running", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artemis-brokers"
		expectedPodCount := 3 // Set your expected number of 'broker' pods

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "application=ex-aao-app"})
		Expect(err).To(BeNil(), "Error getting pods: %v", err)

		// Debugging statements
		fmt.Printf("Retrieved %d pods in namespace %s\n", len(pods.Items), namespace)
		for _, pod := range pods.Items {
			fmt.Printf("Pod Name: %s\n", pod.Name)
			// Add more details as needed
		}

		actualPodCount := len(pods.Items)

		Expect(actualPodCount).To(Equal(expectedPodCount), "Expected %d 'broker' pods, but found %d", expectedPodCount, actualPodCount)
	})
})
