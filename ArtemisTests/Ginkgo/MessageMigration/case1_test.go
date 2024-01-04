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

		// Loop through the found pods to find the specific message
		for _, pod := range pods.Items {
			podName := pod.Name
			receiver, err = session.NewReceiver(
				amqp.LinkSourceAddress(queueName),
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			fmt.Printf("Receiver created for pod '%s'.\n", podName)

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
