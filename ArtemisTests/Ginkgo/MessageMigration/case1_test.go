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
	var clients []*amqp.Client
	var sessions []*amqp.Session
	var senders []*amqp.Sender
	var receivers map[string]*amqp.Receiver
	var ctx context.Context
	var err error

	var kubeClient *kubernetes.Clientset
	var namespace string

	ginkgo.BeforeEach(func() {
		// (Your existing setup code remains unchanged)
	})

	ginkgo.It("should send, delete, and check messages", func() {
		// (Your existing setup code remains unchanged)

		queueName := "doei"
		messageText := "doei"

		// Step 1: Establish connections to multiple brokers and create senders
		brokers := []string{"ex-aao-ss-0", "ex-aao-ss-1", "ex-aao-ss-2"}
		for _, broker := range brokers {
			client, session, sender, err := createConnectionAndSender(broker, queueName)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			clients = append(clients, client)
			sessions = append(sessions, session)
			senders = append(senders, sender)
		}

		// Step 2: Wait for a short duration
		time.Sleep(20 * time.Second)

		// Step 3: Determine which broker received the message
		for _, receiver := range receivers {
			// (Your existing message receiving logic remains unchanged)
		}

		// Step 4: Delete the broker that received the message
		deleteBroker := receivedBroker
		DeleteBroker(deleteBroker, "activemq-artemis-brokers", kubeClient)

		// Step 5: Wait for 60 seconds
		time.Sleep(60 * time.Second)

		// Step 6: Check in the remaining brokers where the message is found
		for _, receiver := range receivers {
			// (Your existing message receiving logic remains unchanged)
		}

		// Step 7: Print a message based on the search status
		if messageFound {
			fmt.Println("Message search completed. Message found.")
		} else {
			fmt.Println("Message search completed. Message not found.")
		}
	})

	ginkgo.AfterEach(func() {
		// (Your existing cleanup code remains unchanged)
		// Close all resources
		for _, sender := range senders {
			sender.Close(ctx)
		}
		for _, receiver := range receivers {
			receiver.Close(ctx)
		}
		for _, session := range sessions {
			session.Close(ctx)
		}
		for _, client := range clients {
			client.Close()
		}
	})
})

func createConnectionAndSender(broker, queueName string) (*amqp.Client, *amqp.Session, *amqp.Sender, error) {
	client, err := amqp.Dial(
		"amqp://" + broker + ".activemq-artemis-brokers.svc.cluster.local:61619",
		amqp.ConnSASLPlain("cgi", "cgi"),
		amqp.ConnIdleTimeout(30*time.Second),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, nil, err
	}

	sender, err := session.NewSender(
		amqp.LinkTargetAddress(queueName),
		amqp.LinkSourceAddress(queueName),
	)
	if err != nil {
		session.Close(context.Background())
		client.Close()
		return nil, nil, nil, err
	}

	// Add the receiver to the map (you can adapt this logic based on your actual code structure)
	receivers[broker], err = session.NewReceiver(
		amqp.LinkSourceAddress(queueName),
	)
	if err != nil {
		sender.Close(context.Background())
		session.Close(context.Background())
		client.Close()
		return nil, nil, nil, err
	}

	return client, session, sender, nil
}
