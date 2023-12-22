package NodeAffinity_test

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("ActiveMQ Artemis Node Affinity Test", func() {
	var kubeClient *kubernetes.Clientset
	var namespace string

	ginkgo.BeforeEach(func() {
		var err error
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		namespace = "activemq-artemis-brokers"
	})

	ginkgo.It("Should have ActiveMQArtemis pods on different nodes", func() {
		// Get the list of broker pods in the namespace
		pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "ex-aao-broker", // Update with the actual label selector for broker pods
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error getting broker pod list")

		// Print debugging information
		fmt.Printf("Namespace: %s\n", namespace)
		fmt.Printf("Found %d broker pods\n", len(pods.Items))

		// Check if broker pods are on different nodes
		nodes := make(map[string]struct{})
		for _, pod := range pods.Items {
			nodeName := pod.Spec.NodeName
			_, exists := nodes[nodeName]
			gomega.Expect(exists).To(gomega.BeFalse(), fmt.Sprintf("Broker Pod %s is on the same node as another pod", pod.Name))
			nodes[nodeName] = struct{}{}

			// Print the pod name and associated node
			fmt.Printf("Broker Pod Name: %s, Node: %s\n", pod.Name, nodeName)
		}

		// Confirm that broker pods are on different nodes
		fmt.Println("All ActiveMQArtemis broker pods are on different nodes.")
	})
})
