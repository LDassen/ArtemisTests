package AutoCreationQueue_test

import (
    "context"
    "github.com/onsi/ginkgo/v2"
    "github.com/onsi/gomega"
    "pack.ag/amqp" // AMQP library for Go
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strings"
)

var _ = ginkgo.Describe("Artemis Queue Test with AMQP", func() {
    var client *amqp.Client
    var session *amqp.Session
    var sender *amqp.Sender
    var receiver *amqp.Receiver
    var ctx context.Context
    var err error

    ginkgo.BeforeEach(func() {
        ctx = context.Background()
        client, err = amqp.Dial("amqp://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619", amqp.ConnSASLPlain("cgi", "cgi")) // Replace with actual credentials and Artemis server address
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        session, err = client.NewSession()
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
    })

    ginkgo.It("should send and receive a message in a queue", func() {
        queueName := "TESTKUBE"
        messageText := "Hello, Artemis!"

        // Create a sender
        sender, err = session.NewSender(
            amqp.LinkTargetAddress(queueName),
        )
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        // Send a message
        err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        // Create a receiver
        receiver, err = session.NewReceiver(
            amqp.LinkSourceAddress(queueName),
        )
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        // Receive a message
        msg, err := receiver.Receive(ctx)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        gomega.Expect(string(msg.GetData())).To(gomega.Equal(messageText))

        // Accept message
        msg.Accept()
    })

    ginkgo.AfterEach(func() {
        if sender != nil {
            sender.Close(ctx)
        }
        if receiver != nil {
            receiver.Close(ctx)
        }
        if session != nil {
            session.Close(ctx)
        }
        if client != nil {
            client.Close()
        }
    })

    // Function to check if a queue exists
func checkQueueExists(queueName string) bool {
    url := "http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:8161/console/jolokia/exec/org.apache.activemq.artemis:broker=\"0.0.0.0\"/listQueues/{\"field\":\"name\",\"operation\":\"CONTAINS\",\"value\":\"" + queueName + "\"}/1/100"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false
    }
    req.Header.Set("Origin", "http://localhost")
    req.SetBasicAuth("cgi", "cgi") // Replace with actual credentials

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return false
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return false
    }

    // Parse the JSON response
    var queues []map[string]interface{}
    if err := json.Unmarshal(body, &queues); err != nil {
        return false
    }

    for _, queue := range queues {
        if name, ok := queue["name"].(string); ok && strings.Contains(name, queueName) {
            return true
        }
    }

    return false
}

    // Add a test case to check for the queue
    ginkgo.Describe("Queue Existence Check", func() {
        ginkgo.It("should verify the existence of the TESTKUBE queue", func() {
            exists := checkQueueExists("TESTKUBE")
            gomega.Expect(exists).To(gomega.BeTrue(), "Queue TESTKUBE should exist")
        })
    })
})
