package SSLConfig_test

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"github.com/IBM/sarama"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kafka SSL Connection", func() {
	Context("when connecting to Kafka over SSL", func() {
		It("should establish a secure connection successfully", func() {
			certPool := x509.NewCertPool()
			caCert, err := ioutil.ReadFile("/var/kafka/ca.crt")
			Expect(err).NotTo(HaveOccurred())
			certPool.AppendCertsFromPEM(caCert)

			cert, err := tls.LoadX509KeyPair("/var/kafka/tls.crt", "/var/kafka/tls.key")
			Expect(err).NotTo(HaveOccurred())

			config := sarama.NewConfig()
			config.Net.TLS.Enable = true
			config.Net.TLS.Config = &tls.Config{
				RootCAs:      certPool,
				Certificates: []tls.Certificate{cert},
			}

			broker := "kafka-brokers-headless.kafka-brokers.svc.cluster.local:9094"
			consumer, err := sarama.NewConsumer([]string{broker}, config)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				if err := consumer.Close(); err != nil {
					log.Println("Error closing consumer:", err)
				}
			}()

			// Perform additional tests to verify the connection
			// For example, you might list topics, consume messages, etc.
			_, err = consumer.Topics()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
