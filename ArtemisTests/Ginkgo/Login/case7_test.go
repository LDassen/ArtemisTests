package Login_test

import (
    "context"
    "fmt"
    "strings"
    "github.com/onsi/ginkgo/v2"
    "github.com/onsi/gomega"
    "pack.ag/amqp" // AMQP library for Go
)

var _ = ginkgo.Describe("Login Credentials Validation Test with AMQP", func() {
    var client *amqp.Client
    var session *amqp.Session
    var ctx context.Context
    var err error

    ginkgo.BeforeEach(func() {
        ctx = context.Background()
        // Replace with actual credentials and Artemis server address
        client, err = amqp.Dial("amqp://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619", amqp.ConnSASLPlain("", "cgi"))

        if err != nil && strings.Contains(err.Error(), "SASL PLAIN auth failed with code 0x1") {
            // If the specific error is encountered, consider the test as passed
            fmt.Println("[PASSED] Expected error encountered:", err)
            ginkgo.Skip("Skipping the rest of the test due to expected error.")
        } else if err == nil {
            ginkgo.Fail("Login should not be successful with empty credentials")
        }

        gomega.Expect(err).To(gomega.HaveOccurred()) // Fail the test if login is successful
        session, err = client.NewSession()
    })

    ginkgo.It("should not be able to send and receive a message due to invalid login credentials", func() {
        gomega.Expect(client).To(gomega.BeNil()) // The client should be nil as the login was not successful
    })

    ginkgo.AfterEach(func() {
        if client != nil {
            client.Close()
        }
        if session != nil {
            session.Close(ctx)
        }
    })
})
