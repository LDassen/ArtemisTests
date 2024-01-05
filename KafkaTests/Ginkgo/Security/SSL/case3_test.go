package SSLConfig_test

import (
	"log"
	"time"
	"github.com/IBM/sarama"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Kafka Simple Connection", func() {
	Context("when sending and receiving messages", func() {
		It("should fail to produce and consume messages on TESTKUBE topic", func() {
			config := sarama.NewConfig()
			config.Producer.Return.Successes = true

			broker := "kafka-brokers-headless.kafka-brokers.svc.cluster.local:9094"

			// Producing a message
			producer, err := sarama.NewSyncProducer([]string{broker}, config)
			if err == nil {
				defer func() {
					if err := producer.Close(); err != nil {
						log.Println("Error closing producer:", err)
					}
				}()

				message := "Hello, TestKube!"
				_, _, err := producer.SendMessage(&sarama.ProducerMessage{
					Topic: "TESTKUBE",
					Value: sarama.StringEncoder(message),
				})

				// Test fails if message is successfully sent
				Expect(err).To(HaveOccurred())
			} else {
				// Test succeeds if producer cannot be created
				Expect(err).To(HaveOccurred())
			}

			// Consuming the message
			consumer, err := sarama.NewConsumer([]string{broker}, config)
			if err == nil {
				defer func() {
					if err := consumer.Close(); err != nil {
						log.Println("Error closing consumer:", err)
					}
				}()

				partitionConsumer, err := consumer.ConsumePartition("TESTKUBE", 0, sarama.OffsetNewest)
				if err == nil {
					defer func() {
						if err := partitionConsumer.Close(); err != nil {
							log.Println("Error closing partition consumer:", err)
						}
					}()

					select {
					case <-partitionConsumer.Messages():
						// Test fails if message is successfully consumed
						Fail("Message was unexpectedly received")
					case <-time.After(5 * time.Second):
						// Test succeeds if no message is received within the timeout
					}
				} else {
					// Test succeeds if consumer cannot be created or consume partition
					Expect(err).To(HaveOccurred())
				}
			} else {
				// Test succeeds if consumer cannot be created
				Expect(err).To(HaveOccurred())
			}
		})
	})
})
