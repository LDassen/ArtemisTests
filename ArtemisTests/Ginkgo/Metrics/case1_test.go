package Metrics_test

import (
    "net/http"
    "io/ioutil"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("Artemis Metrics", func() {
    It("should retrieve metrics successfully", func() {
        url := "http://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:8161/metrics"

        // Perform the HTTP GET request
        resp, err := http.Get(url)
        Expect(err).NotTo(HaveOccurred())
        defer resp.Body.Close()

        // Read the response body
        body, err := ioutil.ReadAll(resp.Body)
        Expect(err).NotTo(HaveOccurred())

        // Convert the body to a string for checking
        bodyString := string(body)

        // Check the HTTP status code and body content
        Expect(resp.StatusCode).To(Equal(http.StatusOK))
        Expect(bodyString).To(ContainSubstring("your_metric_keyword")) // Replace with actual keyword
    })
})
