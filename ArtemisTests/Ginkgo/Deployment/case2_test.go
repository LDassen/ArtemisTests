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
	It("should have one operator pod running in 'activemq-artermis-operator' namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artermis-operator"
		expectedPodPattern := "activemq-artemis-controller-manager-[a-z0-9]+-[a-z0-9]+"

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error getting pods: %v", err)

		var actualPodCount int
		var podName string
		for _, pod := range pods.Items {
			match, _ := regexp.MatchString(expectedPodPattern, pod.Name)
			if match && pod.Status.Phase == "Running" {
				fmt.Printf("Operator Pod Name: %s\n", pod.Name)
				actualPodCount++
				podName = pod.Name
			}
		}

		Expect(actualPodCount).To(Equal(1), "Expected one 'operator-pod' to be running in namespace 'activemq-artermis-operator', but found %d", actualPodCount)
		fmt.Printf("Found Operator Pod: %s\n", podName)
	})
})