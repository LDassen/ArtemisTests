package SSLConfig_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/Shopify/sarama"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
)

var _ = Describe("Kafka SSL Test", func() {
    var headless string      // Single headless address
    var topic string       // Kafka topic
    var tlsConfig *tls.Config

    BeforeEach(func() {
        // Load the CA certificate
        caCert, err := ioutil.ReadFile("/etc/ssl/certs/kafka-bundle.pem")
        Expect(err).NotTo(HaveOccurred())

        caCertPool := x509.NewCertPool()
        caCertPool.AppendCertsFromPEM(caCert)

        tlsConfig = &tls.Config{
            RootCAs: caCertPool,
        }

        // Set to the DNS name of your headless service and port
        headless = "kafka-brokers-headless.kafka-brokers.svc.cluster.local:9094" // Replace with your headless service
        topic = "TESTKUBE"                   // Replace with your Kafka topic
    })

    It("should successfully produce and consume messages with Kafka over SSL", func() {
        config := sarama.NewConfig()
        config.Net.TLS.Enable = true
        config.Net.TLS.Config = tlsConfig
        config.Producer.Return.Successes = true

        // Create a new Sarama producer
        producer, err := sarama.NewSyncProducer([]string{broker}, config)
        Expect(err).NotTo(HaveOccurred())
        defer producer.Close()

        // Produce a message
        _, _, err = producer.SendMessage(&sarama.ProducerMessage{
            Topic: topic,
            Value: sarama.StringEncoder("Hello Kafka over SSL!"),
        })
        Expect(err).NotTo(HaveOccurred())

        // Create a new Sarama consumer
        consumer, err := sarama.NewConsumer([]string{broker}, config)
        Expect(err).NotTo(HaveOccurred())
        defer consumer.Close()

        partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
        Expect(err).NotTo(HaveOccurred())
        defer partitionConsumer.Close()

        // Consume a message
        consumedMessage := <-partitionConsumer.Messages()
        Expect(string(consumedMessage.Value)).To(Equal("Hello Kafka over SSL!"))
    })
})