package AutoCreationQueue_test

import (
    "context"
    "github.com/streadway/amqp"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    v1 "k8s.io/api/core/v1"
)

var _ = Describe("Broker Label Test", func() {
    var (
        clientset *kubernetes.Clientset
        namespace string = "activemq-artemis-brokers"
    )

    BeforeEach(func() {
        // Set up the client
        config, err := rest.InClusterConfig()
        Expect(err).NotTo(HaveOccurred())
        clientset, err = kubernetes.NewForConfig(config)
        Expect(err).NotTo(HaveOccurred())

        // Connect to Artemis using AMQP
        conn, err := amqp.Dial("amqp://cgi:cgi@ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619")
        Expect(err).NotTo(HaveOccurred())
        defer conn.Close()

        ch, err := conn.Channel()
        Expect(err).NotTo(HaveOccurred())
        defer ch.Close()

        body := "Hi, this is a test!"
        err = ch.Publish(
            "",         // exchange
            "TESTKUBE", // routing key (queue name)
            false,      // mandatory
            false,      // immediate
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(body),
            })
        Expect(err).NotTo(HaveOccurred())
    })

    It("should execute a command in one of the pods", func() {
        pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
            LabelSelector: "application=ex-aao-app",
        })
        Expect(err).NotTo(HaveOccurred())

        // Assuming you pick the first pod for simplicity
        podName := pods.Items[0].Name
        cmd := []string{"bin/bash", "-c", "./amq-broker/bin/artemis queue stat --url tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616 --user cgi --password cgi --maxRows 200"}

        req := clientset.CoreV1().Pods(namespace).Exec(podName, &v1.PodExecOptions{
            Command: cmd,
            Stdout:  true,
            Stderr:  true,
        })

        // Execute the command
        // Note: You might need to handle the stream (stdout, stderr) to capture the command output
        _, err = Execute(req, config)
        Expect(err).NotTo(HaveOccurred())

        // Add your logic here to check for the specific number in the output line with 'TEST'
    })
})