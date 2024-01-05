package SSL_test

import (
	"log"
	"time"
	"github.com/IBM/sarama"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kafka SSL Connection", func() {
	Context("when sending and receiving messages over SSL", func() {
		It("should not successfully produce and consume messages on TESTKUBE topic", func() {

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
			partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: "TESTKUBE",
				Value: sarama.StringEncoder(message),
			})
			Expect(err).To(HaveOccurred())

			// Consuming the message
			consumer, err := sarama.NewConsumer([]string{broker}, config)
			Expect(err).To(HaveOccurred())
			defer func() {
				if err := consumer.Close(); err != nil {
					log.Println("Error closing consumer:", err)
				}
			}()

			partitionConsumer, err := consumer.ConsumePartition("TESTKUBE", partition, offset)
			Expect(err).To(HaveOccurred())
			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					log.Println("Error closing partition consumer:", err)
				}
			}()

			select {
			case msg := <-partitionConsumer.Messages():
				Expect(string(msg.Value)).To(Equal(message))
			case <-time.After(5 * time.Second):
				Fail("Timed out waiting for message")
			}
		})
	})
})
