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

var _ = Describe("Check the ActiveMQ Artemis Operator Pod", func() {
	It("should have one operator pod running in 'activemq-artermis-operator' namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artermis-operator"
		labelSelector := "ex-aao"

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		Expect(err).To(BeNil(), "Error getting pods: %v", err)

		Expect(len(pods.Items)).To(BeNumerically(">", 0), "Expected at least one pod with label '%s' to be running in namespace '%s', but found none", labelSelector)

		fmt.Printf("Found at least one Operator Pod with label '%s' in namespace '%s'\n", labelSelector, namespace)
	})
})