package AutoCreationQueue_test

import (
    "context"
    "github.com/onsi/ginkgo/v2"
    "github.com/onsi/gomega"
    "pack.ag/amqp"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "strings"
)

// Function to check if a queue exists
func checkQueueExists(queueName string) bool {
    url := "http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:8161/console/jolokia/exec/org.apache.activemq.artemis:broker=\"0.0.0.0\"/listQueues/{\"field\":\"name\",\"operation\":\"CONTAINS\",\"value\":\"" + queueName + "\"}/1/100"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false
    }
    req.Header.Set("Origin", url)
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

    ginkgo.It("should send a message and then check if the queue exists", func() {
        queueName := "TESTKUBE"
        messageText := "Hello, Artemis!"

        // Create a sender and send a message
        sender, err = session.NewSender(amqp.LinkTargetAddress(queueName))
        gomega.Expect(err).NotTo(gomega.HaveOccurred())
        err = sender.Send(ctx, amqp.NewMessage([]byte(messageText)))
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        // Check if the queue exists after sending the message
        exists := checkQueueExists(queueName)
        gomega.Expect(exists).To(gomega.BeTrue(), "Queue TESTKUBE should be created after message is sent")
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
})
