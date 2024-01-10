package Deployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

var _ = Describe("Check the ActiveMQ Artemis Broker Pods", func() {
	It("should have the correct 'ex-aao-ss-' prefixed pods running", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artemis-brokers"
		expectedPodPrefix0 := "ex-aao-ss-0"
		expectedPodPrefix1 := "ex-aao-ss-1"
		expectedPodPrefix2 := "ex-aao-ss-2"

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error getting pods: %v", err)

		var actualPodCount int
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, expectedPodPrefix0) && pod.Status.Phase == "Running" {
				fmt.Printf("Pod Name0: %s\n", pod.Name)
				actualPodCount++
			}
			if strings.HasPrefix(pod.Name, expectedPodPrefix1) && pod.Status.Phase == "Running" {
				fmt.Printf("Pod Name1: %s\n", pod.Name)
				actualPodCount++
			}
			if strings.HasPrefix(pod.Name, expectedPodPrefix2) && pod.Status.Phase == "Running" {
				fmt.Printf("Pod Name2: %s\n", pod.Name)
				actualPodCount++
			}
		}

		// Set your expected number of pods here
		expectedPodCount := 3 
		Expect(actualPodCount).To(Equal(expectedPodCount), "Expected %d 'ex-aao-ss-' pods, but found %d", expectedPodCount, actualPodCount)
	})
})
