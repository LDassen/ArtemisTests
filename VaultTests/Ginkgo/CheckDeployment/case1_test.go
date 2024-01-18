package CheckDeployment_test

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

var _ = Describe("Vault instances", func() {
	It("should have the correct vault instances instances running", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "dev-vault-1"
		expectedPodPrefix0 := "dev-vault-1-0"
		expectedPodPrefix1 := "dev-vault-1-1"
		expectedPodPrefix2 := "dev-vault-1-2"
		expectedPodPrefix3 := "dev-vault-1-3"
		expectedPodPrefix4 := "dev-vault-1-4"

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
			if strings.HasPrefix(pod.Name, expectedPodPrefix3) && pod.Status.Phase == "Running" {
				fmt.Printf("Pod Name3: %s\n", pod.Name)
				
				actualPodCount++
			}
			if strings.HasPrefix(pod.Name, expectedPodPrefix4) && pod.Status.Phase == "Running" {
				fmt.Printf("Pod Name4: %s\n", pod.Name)
				
				actualPodCount++
			}
		}

		// Set your expected number of vault-instances pods here
		expectedPodCount := 5
		Expect(actualPodCount).To(Equal(expectedPodCount), "Expected %d vault instances, but found %d", expectedPodCount, actualPodCount)
	})
})
