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
		queueName := "zoz4"
		messageText := "zozmessage4"
	
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
		time.Sleep(50 * time.Second)
	
		// Step 3: Determine which broker received the message
		var receivedBroker string
		messageFound := false
		var cancel context.CancelFunc // Declare cancel variable
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
	
				// Set the flag to true
				messageFound = true
	
				// Exit the loop as the message is found
				break
			}
	
			// Close the receiver
			receiver.Close(ctx)
		}
	
		// Close the receiver after finishing the loop
		if receiver != nil {
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
		time.Sleep(50 * time.Second)
	
		// Step 6: Re-determine which broker currently has the message
		receiver, err = session.NewReceiver(
			amqp.LinkSourceAddress(queueName),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	
		ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
	
		msg, err := receiver.Receive(ctx)
		if err == context.DeadlineExceeded {
			// If timeout is reached, consider the message not found
			fmt.Println("Timeout exceeded while searching for the message.")
		} else if err != nil {
			// Handle other errors
			fmt.Printf("Error receiving message: %v\n", err)
		} else {
			// Check if the received message matches the specific message
			if string(msg.GetData()) == messageText {
				// Print where the message was found
				fmt.Printf("Message found in broker '%s'.\n", receivedBroker)
	
				// Accept the message
				msg.Accept()
	
				// Set the flag to true
				messageFound = true
			}
	
			// Close the receiver
			receiver.Close(ctx)
		}
	
		// Step 7: Print a message based on the search status
		if messageFound {
			fmt.Println("Message search completed. Message found.")
		} else {
			fmt.Println("Message search completed. Message not found.")
		}
	})	

	ginkgo.AfterEach(func() {
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
