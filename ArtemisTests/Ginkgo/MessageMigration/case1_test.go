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

		// Establish connection to the Artemis broker
		client, err = amqp.Dial(
			"amqp://ex-aao-ss-2.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
			amqp.ConnSASLPlain("cgi", "cgi"),
			amqp.ConnIdleTimeout(30*time.Second),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		fmt.Println("Connected successfully.")

		// Create a session
		session, err = client.NewSession()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		fmt.Println("Session created successfully.")

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

		// Create a sender and send a message to the specific queue in ex-aao-ss-2 broker
		sender, err = session.NewSender(
			amqp.LinkTargetAddress(queueName),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		fmt.Println("Message sent successfully.")

		// Wait for a short duration
		time.Sleep(10 * time.Second)

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

			// Retry connecting to AMQP after deleting the pod
			retryAttempts := 3
			for attempt := 1; attempt <= retryAttempts; attempt++ {
				// Attempt to connect to AMQP
				client, err = amqp.Dial(
					"amqp://ex-aao-ss-2.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
					amqp.ConnSASLPlain("cgi", "cgi"),
					amqp.ConnIdleTimeout(30*time.Second),
				)

				// Check if the connection was successful
				if err == nil {
					fmt.Println("Connected successfully after retry.")
					break
				}

				// If not successful, wait for some time before retrying
				fmt.Printf("Failed to connect (attempt %d/%d), waiting before retry...\n", attempt, retryAttempts)
				time.Sleep(10 * time.Second)
			}

			// Check if the connection was successful after retrying
			if err != nil {
				gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error during connection after retry: %v", err)
			}

			// Create a new session
			session, err = client.NewSession()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			fmt.Println("Session created successfully after retry.")
		} else {
			fmt.Println("No pods found with label 'application=ex-aao-app'")
		}

		// Loop through the pod names (ex-aao-ss-0, ex-aao-ss-1) to find the specific message
		for _, podName := range []string{"ex-aao-ss-0", "ex-aao-ss-1"} {
			receiver, err = session.NewReceiver(
				amqp.LinkSourceAddress(queueName),
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			fmt.Printf("Receiver created for pod '%s'.\n", podName)

			go func(podName string) {
				// Receive messages from the queue
				msg, err := receiver.Receive(ctx)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				fmt.Printf("Received message from pod '%s': %s\n", podName, string(msg.GetData()))

				// Check if the received message matches the specific message
				gomega.Expect(string(msg.GetData())).To(gomega.Equal(messageText))

				// Accept the message
				msg.Accept()
				fmt.Printf("Message accepted from pod '%s'.\n", podName)

				// Close the receiver
				receiver.Close(ctx)
				fmt.Printf("Receiver closed for pod '%s'.\n", podName)
			}(podName)
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
