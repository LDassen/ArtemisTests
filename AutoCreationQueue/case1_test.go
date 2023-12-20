package AutoCreationQueue_test

import (
    "context"
    "github.com/onsi/ginkgo/v2"
    "github.com/onsi/gomega"
    "pack.ag/amqp"
)

// Function to check if a queue exists
func checkQueueExists(queueName string) bool {
    // ... existing code for checkQueueExists ...
}

var _ = ginkgo.Describe("Artemis Queue Test with AMQP", func() {
    var client *amqp.Client
    var session *amqp.Session
    var sender *amqp.Sender
    var receiver *amqp.Receiver
    var ctx context.Context
    var err error

    ginkgo.BeforeEach(func() {
        ctx = context.Background()
        client, err = amqp.Dial("amqp://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619", amqp.ConnSASLPlain("cgi", "cgi")) // Replace with actual credentials and Artemis server address
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        session, err = client.NewSession()
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    })

    ginkgo.It("should send a message and then check if the queue exists", func() {
        queueName := "TESTKUBE"
        messageText := "Hello, Artemis!"

        // Create a sender and send a message
        sender, err = session.NewSender(amqp.LinkTargetAddress(queueName))
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        // Check if the queue exists after sending the message
        exists := checkQueueExists(queueName)
        gomega.Expect(exists).To(gomega.BeTrue(), "Queue TESTKUBE should be created after message is sent")
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
