package MessageMigration_test

import (
	"context"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"pack.ag/amqp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = ginkgo.Describe("Artemis Queue Specific Broker Test", func() {
	var client *amqp.Client
	var session *amqp.Session
	var sender *amqp.Sender
	var receiver *amqp.Receiver
	var ctx context.Context
	var err error
	var kubeClient *kubernetes.Clientset // Replace with your Kubernetes client initialization

	var specificBrokerAddress string = "amqp://ex-aao-ss-2.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61617"
	var otherBrokerAddresses []string = []string{"amqp://ex-aao-ss-0.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61617", "amqp://ex-aao-ss-1.ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61617"} // Replace with actual addresses of other brokers

	ginkgo.BeforeEach(func() {
		ctx = context.Background()

		fmt.Println("Connecting to the Artemis broker...")
		// Establish connection to the specific Artemis broker
		client, err = amqp.Dial(
			specificBrokerAddress,
			amqp.ConnSASLPlain("cgi", "cgi"),
			amqp.ConnIdleTimeout(30*time.Second),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		session, err = client.NewSession()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("should send a message, delete the broker, and check message in other brokers", func() {
		queueName := "TTTTTTESTKUBE"
		messageText := "Hi, this is a test message"

		// Create a sender
		sender, err = session.NewSender(
			amqp.LinkTargetAddress(queueName),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Send the message
		err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Delete the last broker's pod
		deleteBrokerAddr := otherBrokerAddresses[len(otherBrokerAddresses)-1]
		deletePodName := fmt.Sprintf("ex-aao-%d", len(otherBrokerAddresses)-1)
		deletePodNamespace := "activemq-artemis-brokers"
		deletePropagationPolicy := metav1.DeletePropagationForeground
		deleteOptions := &metav1.DeleteOptions{PropagationPolicy: &deletePropagationPolicy}
		if err := kubeClient.CoreV1().Pods(deletePodNamespace).Delete(ctx, deletePodName, *deleteOptions); err != nil {
			ginkgo.Fail(fmt.Sprintf("Failed to delete pod: %v", err))
			return
		}

		// Check other brokers for the message
		for _, brokerAddress := range otherBrokerAddresses {
			// Connect to another broker
			client, err = amqp.Dial(brokerAddress, amqp.ConnSASLPlain("cgi", "cgi")) // Adjust as necessary
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			session, err = client.NewSession()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Create a receiver
			receiver, err = session.NewReceiver(
				amqp.LinkSourceAddress(queueName),
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Attempt to receive the message
			msg, err := receiver.Receive(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(string(msg.GetData())).To(gomega.Equal(messageText))

			// Clean up
			receiver.Close(ctx)
			session.Close(ctx)
			client.Close()
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
