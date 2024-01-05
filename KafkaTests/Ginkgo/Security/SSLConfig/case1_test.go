package SSLConfig_test

import (
    // ... other imports ...
    "github.com/IBM/sarama"
    "fmt"
    "time"
)

var _ = Describe("Kafka SSL Test", func() {
    // ... (other code remains the same) ...

    It("should successfully produce and consume messages with Kafka over SSL", func() {
        // ... (producer setup remains the same) ...

        // Produce a message with a timeout
        fmt.Println("Attempting to produce a message...")
        select {
        case err = <-func() chan error {
            errChan := make(chan error, 1)
            go func() {
                _, _, err = producer.SendMessage(&sarama.ProducerMessage{
                    Topic: topic,
                    Value: sarama.StringEncoder("Hello Kafka over SSL!"),
                })
                errChan <- err
            }()
            return errChan
        }():
            Expect(err).NotTo(HaveOccurred())
        case <-time.After(30 * time.Second): // Timeout after 30 seconds
            Fail("Timeout while producing message")
        }

        fmt.Println("Message produced, setting up consumer...")

        // ... (consumer setup remains the same) ...

        // Consume a message with a timeout
        fmt.Println("Attempting to consume a message...")
        select {
        case consumedMessage := <-partitionConsumer.Messages():
            Expect(string(consumedMessage.Value)).To(Equal("Hello Kafka over SSL!"))
        case <-time.After(30 * time.Second): // Timeout after 30 seconds
            Fail("Timeout while consuming message")
        }
    })
})
