package MessageMigration_test

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/ginkgo"
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
		config, _ := clientcmd.BuildConfigFromFlags("", "")
		kubeClient, _ = kubernetes.NewForConfig(config)

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
		var connectedBrokerAddr string
		for _, brokerAddr := range brokerAddresses {
			client, err := amqp.Dial(
				fmt.Sprintf("amqp://%s", brokerAddr),
				amqp.ConnSASLPlain("cgi", "cgi"),
				amqp.ConnIdleTimeout(30*time.Second),
			)
			if err == nil {
				fmt.Printf("Connected successfully to %s.\n", brokerAddr)
				connectedBrokerAddr = brokerAddr
				break
			} else {
				fmt.Printf("Error connecting to %s: %v\n", brokerAddr, err)
			}
		}

		// Check if a connection was successfully established
		if client == nil {
			ginkgo.Fail("Failed to establish a connection to any broker.")
			return
		}
		fmt.Printf("Connected to broker: %s.\n", connectedBrokerAddr)

		// Create a session
		session, err := client.NewSession()
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("Failed to create session: %v", err))
			return
		}

		// Use the queue name without specifying the broker
		queueName := "Testkube test-queue"

		// Specify the broker as a prefix in the source address when creating the sender
		sourceAddress := connectedBrokerAddr + "." + queueName
		receiver, err := session.NewReceiver(
			amqp.LinkSourceAddress(sourceAddress),
		)
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("Failed to create receiver: %v", err))
			return
		}
		fmt.Printf("Receiver created for broker '%s'.\n", sourceAddress)

		// Create a sender and send a message to the specific queue on the last broker
		sender, err := session.NewSender(
			amqp.LinkTargetAddress(queueName),
			amqp.LinkSourceAddress(sourceAddress),
		)
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("Failed to create sender: %v", err))
			return
		}
		err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("Failed to send message: %v", err))
			return
		}
		fmt.Printf("Message sent successfully to broker '%s'.\n", connectedBrokerAddr)

		// Wait for a few seconds before deleting the last broker
		time.Sleep(5 * time.Second)

		// Delete the last broker
		deleteBrokerAddr := brokerAddresses[len(brokerAddresses)-1]
		deletePodName := fmt.Sprintf("ex-aao-%d", len(brokerAddresses)-1)
		deletePodNamespace := "activemq-artemis-brokers"
		deletePropagationPolicy := metav1.DeletePropagationForeground
		deleteOptions := &metav1.DeleteOptions{PropagationPolicy: &deletePropagationPolicy}
		if err := kubeClient.CoreV1().Pods(deletePodNamespace).Delete(ctx, deletePodName, *deleteOptions); err != nil {
			ginkgo.Fail(fmt.Sprintf("Failed to delete pod: %v", err))
			return
		}
		fmt.Printf("Broker '%s' deleted successfully.\n", deleteBrokerAddr)

		// Loop through all pods with the label to find the specific message
		for _, pod := range pods.Items {
			// Create a receiver for each remaining pod
			receiver, err := session.NewReceiver(
				amqp.LinkSourceAddress(sourceAddress),
			)
			if err != nil {
				ginkgo.Fail(fmt.Sprintf("Failed to create receiver: %v", err))
				return
			}
			fmt.Printf("Receiver created for broker '%s'.\n", sourceAddress)

			// Loop to attempt receiving messages multiple times
			for i := 0; i < 3; i++ {
				// Receive messages from the queue
				msg, err := receiver.Receive(ctx)
				if err != nil {
					fmt.Printf("Error receiving message from pod '%s': %v\n", pod.Name, err)
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
			ginkgo.Fail(fmt.Sprintf("Message not found in pod '%s' after multiple attempts.", pod.Name))
		}
	})

	ginkgo.AfterEach(func() {
		if sender != nil {
			_ = sender.Close(ctx)
		}
		if session != nil {
			_ = session.Close(ctx)
		}
		if client != nil {
			_ = client.Close()
		}
	})
})
