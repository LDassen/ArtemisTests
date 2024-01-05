package SSLConfig_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "github.com/IBM/sarama"
    "crypto/tls"
    "crypto/x509"
    "io/ioutil"
	// "time"
)

var _ = Describe("Kafka SSL Test", func() {
    var headless string
    var topic string
    var tlsConfig *tls.Config

    BeforeEach(func() {
        // Load the CA certificate
        caCertBytes, err := ioutil.ReadFile("/var/kafka/ca.crt") // CA certificate
        Expect(err).NotTo(HaveOccurred())

        // Create a CA certificate pool
        caCertPool := x509.NewCertPool()
        caCertPool.AppendCertsFromPEM(caCertBytes)

        // Load the client certificate and key
        clientCert, err := tls.LoadX509KeyPair("/var/kafka/tls.crt", "/var/kafka/tls.key")
        Expect(err).NotTo(HaveOccurred())

        // Create TLS configuration
        tlsConfig = &tls.Config{
            RootCAs:      caCertPool,
            Certificates: []tls.Certificate{clientCert},
        }

        headless = "kafka-brokers-headless.kafka-brokers:9094"
        topic = "TESTKUBE"
    })

    It("should successfully produce and consume messages with Kafka over SSL", func() {
        config := sarama.NewConfig()
        config.Net.TLS.Enable = true
        config.Net.TLS.Config = tlsConfig
        config.Producer.Return.Successes = true

        // Create a new Sarama producer
        producer, err := sarama.NewSyncProducer([]string{headless}, config)
        Expect(err).NotTo(HaveOccurred())
        defer producer.Close()

        // Produce a message
        _, _, err = producer.SendMessage(&sarama.ProducerMessage{
            Topic: topic,
            Value: sarama.StringEncoder("Hello Kafka over SSL!"),
        })
        Expect(err).NotTo(HaveOccurred())

        // Create a new Sarama consumer
        consumer, err := sarama.NewConsumer([]string{headless}, config)
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