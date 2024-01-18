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

var _ = Describe("Vault operator pod", func() {
	It("should have the correct 'vault-operator-' prefixed pods running", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "vault-operator"
		expectedPodPrefix := "vault-operator-"

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error getting pods: %v", err)

		var actualPodCount int
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, expectedPodPrefix) && pod.Status.Phase == "Running" {
				fmt.Printf("Pod Name: %s\n", pod.Name)
				actualPodCount++
			}
		}

		// Set your expected number of 'vault-operator-' pods here
		expectedPodCount := 1 
		Expect(actualPodCount).To(Equal(expectedPodCount), "Expected %d 'vault-operator-' pods, but found %d", expectedPodCount, actualPodCount)
	})
})
