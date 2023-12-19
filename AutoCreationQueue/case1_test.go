package AutoCreationQueue_test

import (
    "context"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
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
    })

    It("should check for exactly three brokers with the correct labels", func() {
        pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
            LabelSelector: "application=ex-aao-app",
        })
        Expect(err).NotTo(HaveOccurred())

        // Check for exactly three brokers
        Expect(len(pods.Items)).To(Equal(3), "There should be exactly three brokers")

        for _, pod := range pods.Items {
            // Log the pod names
            GinkgoWriter.Printf("Pod %s has the correct label\n", pod.Name)
        }
    })
})
