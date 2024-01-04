package SSL_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "pack.ag/amqp"
    "context"
    "crypto/tls"
    "crypto/x509"
)

var _ = Describe("Artemis SSL and AMQP Test", func() {

    It("should successfully connect", func() {
        // AMQP communication
        client, err := amqp.Dial(
            "amqps://ex-aao-ssl-0-svc.activemq-artemis-brokers.svc:61617",
            amqp.ConnSASLPlain("cgi", "cgi"),
            amqp.ConnTLSConfig(tlsConfig),
        )
        if err != nil {
            if _, ok := err.(x509.UnknownAuthorityError); ok {
                Skip("Skipping test due to certificate signed by unknown authority")
            }
            Fail(err.Error())
        }
        defer client.Close()

        session, err := client.NewSession()
        Expect(err).NotTo(HaveOccurred())

        // Sending a message
        sender, err := session.NewSender(amqp.LinkTargetAddress("SSL"))
        Expect(err).NotTo(HaveOccurred())
        message := "SSL doesn't work!"
        err = sender.Send(context.Background(), amqp.NewMessage([]byte(message)))
        Expect(err).NotTo(HaveOccurred())

        // Receiving a message
        receiver, err := session.NewReceiver(amqp.LinkSourceAddress("SSL"))
        Expect(err).NotTo(HaveOccurred())
        receivedMsg, err := receiver.Receive(context.Background())
        Expect(err).NotTo(HaveOccurred())

        // Check if the message received is as expected
        receivedMessage := string(receivedMsg.GetData())
        if receivedMessage == message {
            Fail("Test failed because the message was received as expected")
        }

        receivedMsg.Accept() // Acknowledge the message
    })
})