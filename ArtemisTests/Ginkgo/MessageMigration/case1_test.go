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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		fmt.Println("Connecting to the Artemis broker...")
		client, err = amqp.Dial(
			"amqp://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619",
			amqp.ConnSASLPlain("cgi", "cgi"),
			amqp.ConnIdleTimeout(30*time.Second),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		fmt.Println("Connected successfully.")

		// Create a session
		fmt.Println("Creating a session...")
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

		// Ensure the StatefulSet (deployment) exists before proceeding
		statefulSetName := "ex-aao-ss"
		_, err = kubeClient.AppsV1().StatefulSets(namespace).Get(ctx, statefulSetName, metav1.GetOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("should send, delete, and check messages", func() {
		queueName := "stop"
		messageText := "nee"

		// Step 1: Create a sender and send a message to the specific queue in the headless connection
		sender, err = session.NewSender(
			amqp.LinkTargetAddress(queueName),
			amqp.LinkSourceAddress(queueName),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Print a message indicating that the message has been sent to the headless connection
		fmt.Printf("Message sent to headless connection.\n")

		// Step 2: Wait for a short duration
		time.Sleep(60 * time.Second)

		// Step 3: Determine which broker received the message
		var receivedBroker string
		for _, broker := range []string{"ex-aao-ss-0", "ex-aao-ss-1", "ex-aao-ss-2"} {
			// Check if the message is present in the current broker
			receiver, err := session.NewReceiver(
				amqp.LinkSourceAddress(queueName),
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Receive messages from the queue
			msg, err := receiver.Receive(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Check if the received message matches the specific message
			if string(msg.GetData()) == messageText {
				// Store the broker where the message was found
				receivedBroker = broker

				// Print where the message was found
				fmt.Printf("Message found in broker '%s'.\n", receivedBroker)

				// Accept the message
				msg.Accept()

				// Close the receiver
				receiver.Close(ctx)

				// Exit the loop as the message is found
				break
			}

			// Close the receiver
			receiver.Close(ctx)
		}

		// Step 4: Delete the broker that received the message
		deletePodName := receivedBroker
		deletePodNamespace := "activemq-artemis-brokers"
		deletePropagationPolicy := metav1.DeletePropagationForeground
		deleteOptions := &metav1.DeleteOptions{PropagationPolicy: &deletePropagationPolicy}
		err = kubeClient.CoreV1().Pods(deletePodNamespace).Delete(ctx, deletePodName, *deleteOptions)
		gomega.Expect(err).To(gomega.BeNil(), "Error deleting pod: %v", err)
		fmt.Printf("Pod '%s' deleted successfully.\n", deletePodName)

		// Step 5: Print a message indicating the start of the search
		fmt.Println("Searching for the message in other brokers...")
		time.Sleep(120 * time.Second)

		// Step 6: Check in the remaining brokers where the message is found
		for _, broker := range []string{"ex-aao-ss-0", "ex-aao-ss-1", "ex-aao-ss-2"} {
			// Skip the deleted broker
			if broker == receivedBroker {
				continue
			}

			receiver, err = session.NewReceiver(
				amqp.LinkSourceAddress(queueName),
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Receive messages from the queue
			msg, err := receiver.Receive(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Check if the received message matches the specific message
			if string(msg.GetData()) == messageText {
				// Print where the message was found
				fmt.Printf("Message found in broker '%s'.\n", broker)

				// Accept the message
				msg.Accept()

				// Close the receiver
				receiver.Close(ctx)

				// Exit the loop as the message is found
				break
			}

			// Close the receiver
			receiver.Close(ctx)
		}

		// Step 7: Print a message indicating the end of the search
		fmt.Println("Message search completed.")
	})

	ginkgo.AfterEach(func() {
		// Close resources in reverse order of creation
		if sender != nil {
			sender.Close(ctx)
		}
		if receiver != nil {
			receiver.Close(ctx)
		}
		if session != nil {
			session.Close(ctx)
		}
		if client != nil {
			client.Close()
		}
	})
})
