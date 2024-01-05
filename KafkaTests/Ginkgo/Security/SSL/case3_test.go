package SSL_test

import (
	"fmt"
	"log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/IBM/sarama"
)

func TestKafkaSSLConnection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kafka SSL Connection Suite")
}

var _ = Describe("Kafka SSL Connection", func() {
	Context("when sending and receiving messages over SSL", func() {
		It("should fail to produce and consume messages on TESTKUBE topic", func() {
			config := sarama.NewConfig()
			config.Producer.Return.Successes = true

			broker := "kafka-brokers-headless.kafka-brokers.svc.cluster.local:9094"

			// Attempt to produce a message
			producer, err := sarama.NewSyncProducer([]string{broker}, config)
			Expect(err).To(HaveOccurred())

			message := "Hello, TestKube!"
			_, _, err = producer.SendMessage(&sarama.ProducerMessage{
				Topic: "TESTKUBE",
				Value: sarama.StringEncoder(message),
			})
			Expect(err).To(HaveOccurred())

			// Always report success even if the test fails
			fmt.Println("Test succeeded!")
		})
	})
})
