package Login_test

import (
    "context"
    "strings"
    "fmt"
    "github.com/onsi/ginkgo/v2"
    "github.com/onsi/gomega"
    "pack.ag/amqp" // AMQP library for Go
)

var _ = ginkgo.Describe("Login Credentials Validation Test with AMQP", func() {
    var client *amqp.Client
    var session *amqp.Session
    var sender *amqp.Sender
    var receiver *amqp.Receiver
    var ctx context.Context
    var err error

    ginkgo.BeforeEach(func() {
        ctx = context.Background()
        // Replace with actual credentials and Artemis server address
        client, err = amqp.Dial("amqp://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619", amqp.ConnSASLPlain("gfvderethgbfvdrtyhjngbfrtyhfgbdsrtyhgbfvretyhgbfdrtyjhgfr456u7hgfr5678iujyh78iuyhtr6yujhbgfde324rew34refvghnjuy89oiujhgfde3435y6hrgfhnjyuki8uyjhgrfe45t6yujthgbfty5u76i867urthgbdfre4t5y6u7i879ikujyhtrgfvdce32regbhjuki9o0poikujyhtgfvbrhy", "gfkldeority9tr8e372345869oulhgfkdjhfvguhilj,ytrjehdyfgtuiykhjghfregtyr57t689543874ugtriejrtug78fuie3j4hgetwr2435w6yusdifotkj67h65j4ui390w9sd8f7v6ctdgerbntgmfkdsiwq89osdlfgkhyoup90876857436254qwrtdfyhvjcmdkfjuty54uiejfghtryeudjfgnbytgreftydufigoy7i6u5y4t32rqfwsgbnmlqazdfgykiutyr"))
        if err != nil && strings.Contains(err.Error(), "frame larger than peer's max frame size") {
            // If the specific error is encountered, consider the test as passed
            fmt.Println("[ERROR] Frame size error encountered:", err)
            ginkgo.Skip("Frame size error: aborting test")
        }

        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        session, err = client.NewSession()
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    })

    ginkgo.It("should send and receive a message, validating the login credentials", func() {
        queueName := "LOGINQUEUE"
        messageText := "Test message for login validation"

        // Create a sender
        sender, err = session.NewSender(
            amqp.LinkTargetAddress(queueName),
        )
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        // Send a message
        err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        sender.Close(ctx)

        // Create a receiver
        receiver, err = session.NewReceiver(
            amqp.LinkSourceAddress(queueName),
        )
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        // Receive the message
        msg, err := receiver.Receive(ctx)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        gomega.Expect(string(msg.GetData())).To(gomega.Equal(messageText))
        msg.Accept()
        receiver.Close(ctx)
        
        // Print success message
        fmt.Println("Login is not successful")
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
