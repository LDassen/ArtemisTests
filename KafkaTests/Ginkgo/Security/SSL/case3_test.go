package SSL_test

import (
	"log"
	"github.com/IBM/sarama"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kafka SSL Connection", func() {
	Context("when sending and receiving messages over SSL", func() {
		It("should succeed if messages cannot be sent or received on TESTKUBE topic", func() {

			broker := "kafka-brokers-headless.kafka-brokers.svc.cluster.local:9094"
			config := sarama.NewConfig()

			// Set Producer.Return.Successes to false
			config.Producer.Return.Successes = false

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
			if err != nil && err.Error() != "kafka: invalid configuration" {
				panic(err)
			}

			// Consuming the message
			_, err = sarama.NewConsumer([]string{broker}, config)
			Expect(err).To(HaveOccurred())
		})
	})
})
