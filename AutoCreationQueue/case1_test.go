package AutoCreationQueue_test

import (
    "context"
    "fmt"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "io"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/remotecommand"
    // Import Artemis client package (example, adjust based on actual package)
    "github.com/zeus-fyi/zeus/pkg/artemis/client"
)

var _ = Describe("Artemis Queue Message Test", func() {
    var (
        clientset *kubernetes.Clientset
        namespace string = "activemq-artemis-brokers"
        artemisClient *client.Client // Artemis client
    )

    BeforeEach(func() {
        // Set up the Kubernetes client
        config, err := rest.InClusterConfig()
        Expect(err).NotTo(HaveOccurred())
        clientset, err = kubernetes.NewForConfig(config)
        Expect(err).NotTo(HaveOccurred())

        // Set up the Artemis client
        artemisClient = client.NewClient("ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local", "cgi", "cgi")
    })

    It("should send a message to Artemis and check output in pod", func() {
        pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
            LabelSelector: "application=ex-aao-app",
        })
        Expect(err).NotTo(HaveOccurred())

        // Send a message to the Artemis queue
        err = artemisClient.SendMessage("TESTKUBE", "Hi, this is a test!")
        Expect(err).NotTo(HaveOccurred())

        // Exec into one of the pods and check output
        req := clientset.CoreV1().Pods(namespace).GetLogs(pods.Items[0].Name, &v1.PodLogOptions{})
        podLogs, err := req.Stream(context.TODO())
        Expect(err).NotTo(HaveOccurred())
        defer podLogs.Close()

        buf := new(bytes.Buffer)
        _, err = io.Copy(buf, podLogs)
        Expect(err).NotTo(HaveOccurred())
        output := buf.String()

        // Check for specific number in line with 'TEST'
        Expect(output).To(ContainSubstring("TEST"))
        // Additional checks based on output format
    })
})
