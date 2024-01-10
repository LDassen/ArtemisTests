package Deployment_test

import (
	"context"
	"fmt"
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Check the Trust Manager Pod", func() {
	It("should find and verify the 'trust-manager-*' pod in the 'cert-manager' namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "cert-manager"
		expectedPodNamePrefix := "trust-manager-"

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error getting pods: %v", err)

		var trustManagerPodFound bool
		for _, pod := range pods.Items {
			if strings.HasPrefix(pod.Name, expectedPodNamePrefix) && pod.Status.Phase == "Running" {
				fmt.Printf("Trust Manager Pod Name: %s\n", pod.Name)
				trustManagerPodFound = true
				break
			}
		}

		Expect(trustManagerPodFound).To(BeTrue(), "Trust Manager pod not found or not running in namespace '%s'", namespace)
	})
})
