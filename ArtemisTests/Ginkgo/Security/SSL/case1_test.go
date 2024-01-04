package SSL_test

import (
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "pack.ag/amqp" 
)

var _ = Describe("Artemis SSL and AMQP Test", func() {
    // ... (setup and teardown code)

    It("should successfully connect", func() {
        caCert, err := ioutil.ReadFile("/etc/ssl/certs/ca-bundle.crt")
        Expect(err).NotTo(HaveOccurred())

        caCertPool := x509.NewCertPool()
        caCertPool.AppendCertsFromPEM(caCert)

        config := &tls.Config{
            RootCAs: caCertPool,
        }

        conn, err := tls.Dial("tcp", "<ARTEMIS_HOST>:<ARTEMIS_PORT>", config)
        Expect(err).NotTo(HaveOccurred())
        defer conn.Close()

        // AMQP communication
        client, err := amqp.Dial("amqps://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61617", amqp.ConnTLSConfig(config))
        Expect(err).NotTo(HaveOccurred())
        defer client.Close()
        
        session, err := client.NewSession()
        Expect(err).NotTo(HaveOccurred())
        
        // Sending a message
        sender, err := session.NewSender(amqp.LinkTargetAddress("SSL"))
        Expect(err).NotTo(HaveOccurred())
        message := "SSL works!"
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
