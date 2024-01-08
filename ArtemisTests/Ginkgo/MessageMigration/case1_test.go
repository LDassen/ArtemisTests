package MessageMigration_test

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"pack.ag/amqp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = ginkgo.Describe("MessageMigration Test", func() {
	var client *amqp.Client
	var session *amqp.Session
	var sender *amqp.Sender
	var receiver *amqp.Receiver
	var ctx context.Context
	var kubeClient *kubernetes.Clientset
	var namespace string
	var pods *v1.PodList
	var messageText = "specialtext"

	ginkgo.BeforeEach(func() {
		ctx = context.Background()

		// Initialize Kubernetes client with in-cluster config
		config, err := clientcmd.BuildConfigFromFlags("", "")
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Set the namespace
		namespace = "activemq-artemis-brokers"
	})

	ginkgo.It("should send, delete, and check messages", func() {
		// List of broker addresses
		brokerAddresses := []string{
			"ex-aao-ss-0.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
			"ex-aao-ss-1.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
			"ex-aao-ss-2.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
		}

		// Try connecting to each broker until a successful connection is established
		for _, brokerAddr := range brokerAddresses {
			client, err = amqp.Dial(
				fmt.Sprintf("amqp://%s", brokerAddr),
				amqp.ConnSASLPlain("cgi", "cgi"),
				amqp.ConnIdleTimeout(30*time.Second),
			)
			if err == nil {
				fmt.Printf("Connected successfully to %s.\n", brokerAddr)
				break
			} else {
				fmt.Printf("Error connecting to %s: %v\n", brokerAddr, err)
			}
		}

		// Check if a connection was successfully established
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to establish a connection to any broker.")
		fmt.Println("Session created successfully.")

		// Create a session
		session, err = client.NewSession()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Use the queue name without specifying the broker
		queueName := "Testkube test-queue"

		// Specify the broker as a prefix in the source address when creating the sender
		sourceAddress := "ex-aao-ss-2." + queueName
		receiver, err = session.NewReceiver(
			amqp.LinkSourceAddress(sourceAddress),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		fmt.Printf("Receiver created for broker '%s'.\n", sourceAddress)

		// Create a sender and send a message to the specific queue on ex-aao-ss-2
		sender, err = session.NewSender(
			amqp.LinkTargetAddress(queueName),
			amqp.LinkSourceAddress(sourceAddress),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		fmt.Printf("Message sent successfully to broker '%s'.\n", sourceAddress)

		// Wait for the last pod deletion to complete
		podDeleted := false
		timeout := time.After(2 * time.Minute)
		for {
			select {
			case <-timeout:
				gomega.Expect(podDeleted).To(gomega.BeTrue(), "Timeout waiting for pod deletion.")
			default:
				// List pods with the label selector "application=ex-aao-app"
				pods, err := kubeClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: "application=ex-aao-app"})
				gomega.Expect(err).To(gomega.BeNil(), "Error getting pods: %v", err)

				// Check if there are any pods with the specified label
				if len(pods.Items) == 0 {
					podDeleted = true
					fmt.Println("Last pod deleted successfully.")
					break
				} else {
					fmt.Printf("Pods still present: %v\n", pods.Items)
				}

				time.Sleep(15 * time.Second)
			}
		}

		// Wait for a few seconds before deleting the pod associated with ex-aao-2
		time.Sleep(5 * time.Second)

		// Delete the pod associated with ex-aao-2
		deletePodName := "ex-aao-2"                        // Replace with the actual pod name
		deletePodNamespace := "activemq-artemis-brokers"   // Replace with the actual namespace
		deletePropagationPolicy := metav1.DeletePropagationForeground
		deleteOptions := &metav1.DeleteOptions{PropagationPolicy: &deletePropagationPolicy}
		err = kubeClient.CoreV1().Pods(deletePodNamespace).Delete(ctx, deletePodName, *deleteOptions)
		gomega.Expect(err).To(gomega.BeNil(), "Error deleting pod: %v", err)
		fmt.Printf("Pod '%s' deleted successfully.\n", deletePodName)

		// Loop through all pods with the label to find the specific message
		for _, pod := range pods.Items {
			// Create a receiver for each remaining pod
			receiver, err = session.NewReceiver(
				amqp.LinkSourceAddress(sourceAddress),
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			fmt.Printf("Receiver created for broker '%s'.\n", sourceAddress)

			// Loop to attempt receiving messages multiple times
			for i := 0; i < 3; i++ {
				// Receive messages from the queue
				msg, err := receiver.Receive(ctx)
				if err != nil {
					fmt.Printf("Error receiving message from pod '%s': %v\n", pod.Name, err)
					gomega.Expect(err).NotTo(gomega.HaveOccurred())
					// Close the receiver if there was an error
					receiver.Close(ctx)
					fmt.Printf("Receiver closed for pod '%s'.\n", pod.Name)
					time.Sleep(5 * time.Second) // Wait before retrying
					continue
				}

				fmt.Printf("Received message from pod '%s': %s\n", pod.Name, string(msg.GetData()))

				// Check if the received message matches the specific message
				if string(msg.GetData()) == messageText {
					// Accept the message
					msg.Accept()
					fmt.Printf("Message found in pod '%s'.\n", pod.Name)
					return // Break out of the loop after finding the message in one pod
				}

				// Close the receiver
				receiver.Close(ctx)
				fmt.Printf("Receiver closed for pod '%s'.\n", pod.Name)
			}

			// If the loop completes without finding the message, fail the test
			ginkgo.Fail("Message not found in pod '%s' after multiple attempts.", pod.Name)
		}
	})

	ginkgo.AfterEach(func() {
		if sender != nil {
			sender.Close(ctx)
		}
		if session != nil {
			session.Close(ctx)
		}
		if client != nil {
			client.Close()
		}
	})
})
