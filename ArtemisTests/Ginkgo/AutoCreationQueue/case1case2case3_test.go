package AutoCreationQueue_test

import (
    "context"
    "github.com/onsi/ginkgo/v2"
    "github.com/onsi/gomega"
    "pack.ag/amqp" // AMQP library for Go
	"time"
)

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

    ginkgo.It("should send and receive two messages in a queue", func() {
        queueName := "TESTKUBE"
        firstMessageText := "Hi, this is message 1"
        secondMessageText := "Hi, this is message 2"
    
        // Create a sender
        sender, err = session.NewSender(
            amqp.LinkTargetAddress(queueName),
        )
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    
        // Send the first message (case 1)
        err = sender.Send(ctx, amqp.NewMessage([]byte(firstMessageText)))
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    
        // Wait for 10 minutes after sending the first message (case 2)
        time.Sleep(1 * time.Minute)

        // Create a receiver for the first message
        receiver, err = session.NewReceiver(
            amqp.LinkSourceAddress(queueName),
        )
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    
        // Receive the first message
        msg, err := receiver.Receive(ctx)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        gomega.Expect(string(msg.GetData())).To(gomega.Equal(firstMessageText))
    
        // Accept the first message
        msg.Accept()
    
        // Close the first receiver to clean up before receiving the second message
        receiver.Close(ctx)
    
        // Send the second message (case 3)
        err = sender.Send(ctx, amqp.NewMessage([]byte(secondMessageText)))
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    
        // Create a receiver for the second message
        receiver, err = session.NewReceiver(
            amqp.LinkSourceAddress(queueName),
        )
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    
        // Receive the second message
        msg, err = receiver.Receive(ctx)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        gomega.Expect(string(msg.GetData())).To(gomega.Equal(secondMessageText))
    
        // Accept the second message
        msg.Accept()
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
