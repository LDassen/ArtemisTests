package SSL_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "pack.ag/amqp"
    "context"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
)

var _ = Describe("Artemis SSL and AMQP Test", func() {
    var config *tls.Config

    BeforeEach(func() {
        // Load the CA certificate
        caCert, err := ioutil.ReadFile("/etc/ssl/certs/kafka-bundle.pem") // Replace with actual path to your CA cert
        Expect(err).NotTo(HaveOccurred())

        // Create a CA certificate pool and add cert to it
        caCertPool := x509.NewCertPool()
        caCertPool.AppendCertsFromPEM(caCert)
        })
    })

    It("should successfully connect", func() {
        // AMQP communication
        client, err := amqp.Dial("amqps://ex-aao-ssl-0-svc.activemq-artemis-brokers.svc:61617", amqp.ConnSASLPlain("cgi", "cgi"))
        Expect(err).NotTo(HaveOccurred())
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
        Expect(string(receivedMsg.GetData())).To(Equal(message))

        receivedMsg.Accept() // Acknowledge the message
    })
})
