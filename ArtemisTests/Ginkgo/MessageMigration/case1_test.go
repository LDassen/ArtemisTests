package MessageMigration_test

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"pack.ag/amqp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = ginkgo.Describe("MessageMigration Test", func() {
	var client *amqp.Client
	var session *amqp.Session
	var sender *amqp.Sender
	var receiver *amqp.Receiver
	var ctx context.Context
	var err error

	var kubeClient *kubernetes.Clientset
	var namespace string

	ginkgo.BeforeEach(func() {
		ctx = context.Background()

		// List of broker addresses
		brokerAddresses := []string{
			"ex-aao-ss-0.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
			"ex-aao-ss-1.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
			"ex-aao-ss-2.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
			// Add more addresses if needed
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

		// Initialize Kubernetes client with in-cluster config
		config, err := clientcmd.BuildConfigFromFlags("", "")
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Set the namespace
		namespace = "activemq-artemis-brokers"
	})

	ginkgo.It("should send, delete, and check messages", func() {
		queueName := "SpecificQueue"
		messageText := "Hello, this is a test message"

		// Create a sender and send a message to the specific queue in the active broker
		sender, err = session.NewSender(
			amqp.LinkTargetAddress(queueName),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		fmt.Println("Message sent successfully.")

		// Wait for a short duration
		time.Sleep(1 * time.Second)

		// List pods with the label selector "application=ex-aao-app"
		pods, err := kubeClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: "application=ex-aao-app"})
		gomega.Expect(err).To(gomega.BeNil(), "Error getting pods: %v", err)

		// Check if there are any pods with the specified label
		if len(pods.Items) > 0 {
			// Get the name of the last pod in the list
			lastPodName := pods.Items[len(pods.Items)-1].Name

			// Delete the last pod
			err := kubeClient.CoreV1().Pods(namespace).Delete(ctx, lastPodName, metav1.DeleteOptions{})
			gomega.Expect(err).To(gomega.BeNil(), "Error deleting pod: %v", err)
			fmt.Printf("Pod '%s' deleted successfully.\n", lastPodName)

			// Wait for the deletion to propagate
			time.Sleep(30 * time.Second)
		} else {
			fmt.Println("No pods found with label 'application=ex-aao-app'")
		}

		// Flag to determine if the message is found
		messageFound := false

		// Loop through all pods with the label to find the specific message
		for _, pod := range pods.Items {
			receiver, err = session.NewReceiver(
				amqp.LinkSourceAddress(queueName),
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			fmt.Printf("Receiver created for pod '%s'.\n", pod.Name)

			// Receive messages from the queue
			msg, err := receiver.Receive(ctx)
			if err != nil {
				fmt.Printf("Error receiving message from pod '%s': %v\n", pod.Name, err)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				// Close the receiver if there was an error
				receiver.Close(ctx)
				fmt.Printf("Receiver closed for pod '%s'.\n", pod.Name)
				continue
			}

			fmt.Printf("Received message from pod '%s': %s\n", pod.Name, string(msg.GetData()))

			// Check if the received message matches the specific message
			if string(msg.GetData()) == messageText {
				// Accept the message
				msg.Accept()
				fmt.Printf("Message foud in pod '%s'.\n", pod.Name)
				// Set the flag to true
				messageFound = true
			}

			// Close the receiver
			receiver.Close(ctx)
			fmt.Printf("Receiver closed for pod '%s'.\n", pod.Name)

			// Break out of the loop after finding the message in one pod
			if messageFound {
				break
			}
		}

		// Check if the message was found in any pod
		gomega.Expect(messageFound).To(gomega.BeTrue(), "Message not found in any pod.")
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
