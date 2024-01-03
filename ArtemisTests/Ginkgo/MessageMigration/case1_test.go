package MessageMigration_test

import (
	"context"
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
		client, err = amqp.Dial(
			"amqp://ex-aao-ss-2.activemq-artemis-brokers.svc.cluster.local:61616",
			amqp.ConnSASLPlain("cgi", "cgi"),
			amqp.ConnIdleTimeout(30*time.Second),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
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

		// Create a sender and send a message to the specific queue in ex-aao-ss-2 broker
		sender, err = session.NewSender(
			amqp.LinkTargetAddress(queueName),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for a short duration
		time.Sleep(1 * time.Second)

		// Delete the ex-aao-ss-2 broker
		err = kubeClient.AppsV1().Deployments(namespace).Delete(ctx, "ex-aao-ss-2", metav1.DeleteOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for the deletion to propagate
		time.Sleep(30 * time.Second)

		// Check queues of ex-aao-ss-0 and ex-aao-ss-1 to find the specific message
		for _, broker := range []string{"ex-aao-ss-0", "ex-aao-ss-1"} {
			receiver, err = session.NewReceiver(
				amqp.LinkSourceAddress(queueName),
			)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Receive messages from the queue
			msg, err := receiver.Receive(ctx)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Check if the received message matches the specific message
			gomega.Expect(string(msg.GetData())).To(gomega.Equal(messageText))

			// Accept the message
			msg.Accept()

			// Close the receiver
			receiver.Close(ctx)
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
