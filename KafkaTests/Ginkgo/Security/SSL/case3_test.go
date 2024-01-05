package SSL_test

import (
	"log"
	"github.com/IBM/sarama"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kafka SSL Connection", func() {
	Context("when sending and receiving messages over SSL", func() {
		It("should pass if messages cannot be sent or received on TESTKUBE topic", func() {

			broker := "kafka-brokers-headless.kafka-brokers.svc.cluster.local:9094"
			config := sarama.NewConfig()

			// Producing a message
			producer, err := sarama.NewSyncProducer([]string{broker}, config)
			Expect(err).NotTo(HaveOccurred())
			defer func() {
				if err := producer.Close(); err != nil {
					log.Println("Error closing producer:", err)
				}
			}()

			message := "Hello, TestKube!"
			_, _, err = producer.SendMessage(&sarama.ProducerMessage{
				Topic: "TESTKUBE",
				Value: sarama.StringEncoder(message),
			})
			if err == nil {
				Fail("Message should not be sent successfully")
			}

			// Consuming the message
			_, err = sarama.NewConsumer([]string{broker}, config)
			if err == nil {
				Fail("Consumer should not be created successfully")
			}
		})
	})
})
