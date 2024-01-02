package MessageMigration_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ = ginkgo.Describe("Message Migration", func() {
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

	ginkgo.Context("Happy Cases", func() {
		ginkgo.It("[case_1] Delete a broker pod with queues and messages on it", func() {
			// Step 1: Get the list of broker pods
			brokerPods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: "ActiveMQArtemis=ex-aao,application=ex-aao-app", // Replace with the actual label for broker pods
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error getting broker pod list")

			// Step 2: Delete one of the broker pods
			if len(brokerPods.Items) > 0 {
				podNameToDelete := brokerPods.Items[0].Name
				err := kubeClient.CoreV1().Pods(namespace).Delete(context.TODO(), podNameToDelete, metav1.DeleteOptions{})
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error deleting broker pod")

				// Step 3: Wait for the drainer log message in the deleted pod
				gomega.Eventually(func() bool {
					podLogs, err := kubeClient.CoreV1().Pods(namespace).GetLogs(podNameToDelete, &v1.PodLogOptions{}).DoRaw(context.TODO())
					if err == nil {
						// Verify in the logs that the drainer pod started
						return strings.Contains(string(podLogs), "Drainer pod started")
					}
					return false
				}, time.Minute, time.Second).Should(gomega.BeTrue(), "Drainer pod not started in logs")

				// Optional: Print debugging information
				fmt.Printf("Test [case_1] completed successfully. Deleted broker pod: %s\n", podNameToDelete)
			} else {
				// Print a message if there are no broker pods to delete
				fmt.Println("No broker pods found to delete")
			}
		})
	})

	ginkgo.AfterEach(func() {
		// Additional cleanup or verification steps after each test
	})
})
