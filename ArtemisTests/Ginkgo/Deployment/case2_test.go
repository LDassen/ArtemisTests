package Deployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
)

var _ = Describe("Check the ActiveMQ Artemis Operator Pod", func() {
	It("should have operator pod running in 'activemq-artermis-operator' namespace with label 'activemq-artemis-controller-manager'", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artermis-operator"
		expectedPodPattern := "activemq-artemis-controller-manager-[a-z0-9]+-[a-z0-9]+"

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error getting pods: %v", err)

		var actualPodCount int
		for _, pod := range pods.Items {
			match, _ := regexp.MatchString(expectedPodPattern, pod.Name)
			if match && pod.Status.Phase == "Running" {
				controllerManagerLabel, found := pod.Labels["activemq-artemis-controller-manager"]
				if found {
					fmt.Printf("Operator Pod Name: %s\n", pod.Name)
					fmt.Printf("Controller Manager Label: %s\n", controllerManagerLabel)
					actualPodCount++
				}
			}
		}

		// Print additional information in case of failure
		if actualPodCount != 1 {
			fmt.Printf("Actual Pod Count: %d\n", actualPodCount)
			for _, pod := range pods.Items {
				fmt.Printf("Pod Name: %s, Phase: %s\n", pod.Name, pod.Status.Phase)
			}
		}

		// Set your expected number of operator pods here
		expectedPodCount := 1
		Expect(actualPodCount).To(Equal(expectedPodCount), "Expected %d 'operator-pod' pod with label 'activemq-artemis-controller-manager', but found %d", expectedPodCount, actualPodCount)
	})
})
