package MessageMigration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	//v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"pack.ag/amqp"
)

var _ = ginkgo.Describe("MessageMigration Test", func() {
	var client *amqp.Client
	var session *amqp.Session
	var sender *amqp.Sender
	var ctx context.Context
	var err error

	var kubeClient *kubernetes.Clientset
	var podName string
	var namespace string

	ginkgo.BeforeEach(func() {
		ctx = context.Background()

		// Establish connection to the Artemis broker
		client, err = amqp.Dial("amqp://ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61616", amqp.ConnSASLPlain("cgi", "cgi"))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		session, err = client.NewSession()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Initialize Kubernetes client
		kubeConfigPath := "/path/to/your/kube/config"
		config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Set your pod name and namespace
		podName = "ex-aao-ss-2"
		namespace = "activemq-artemis-brokers"
	})

	ginkgo.It("should send messages to a specific queue in a pod", func() {
		queueName := "Sabrine"
		messageText := "Hello, this is a test message"

		// Get the pod (this step is optional and can be removed if not used)
		_, err := kubeClient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Create a sender
		sender, err = session.NewSender(
			amqp.LinkTargetAddress(queueName),
		)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Send messages to the specific queue in the pod
		for i := 0; i < 3; i++ {
			err = sender.Send(ctx, amqp.NewMessage([]byte(fmt.Sprintf("%s %d", messageText, i+1))))
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Wait for a short duration between messages
			time.Sleep(1 * time.Second)
		}
	})

	ginkgo.AfterEach(func() {
		if sender != nil {
			sender.Close(ctx)
		}
		if session != nil {
			session.Close(ctx)
		}
		if client != nil {
			client.Close()
		}
	})
})

//func TestQueueSender(t *testing.T) {
//	gomega.RegisterFailHandler(ginkgo.Fail)
//	ginkgo.RunSpecs(t, "MessageMigration Test Suite")
//}
