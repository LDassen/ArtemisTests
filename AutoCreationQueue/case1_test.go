package AutoCreationQueue_test

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/exec"
	"k8s.io/client-go/util/wait/remote"
	"k8s.io/client-go/util/wait/transport"
	"k8s.io/client-go/util/wait/transport/spdy"
)

var _ = ginkgo.Describe("Artemis Test", func() {
	var (
		kubeClient *kubernetes.Clientset
		namespace  = "activemq-artemis-brokers"
		ctx        = context.TODO()
	)

	ginkgo.BeforeSuite(func() {
		// Initialize Kubernetes client
		config, err := rest.InClusterConfig()
		if err != nil {
			home := homedir.HomeDir()
			kubeconfig := filepath.Join(home, ".kube", "config")
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		}
		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("should test Artemis functionality", func() {
		ginkgo.By("Searching for Artemis pods in the namespace")
		podList, err := kubeClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(podList.Items).NotTo(gomega.BeEmpty())

		podName := "ex-aao-ss-0"

		ginkgo.By("Executing command in Artemis pod")
		execCommandInPod(kubeClient, namespace, podName,
			"./amq-broker/bin/artemis producer --user cgi --password cgi --url tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616 --message-count 100 --user cgi --password cgi --maxRows 200")

		time.Sleep(5 * time.Second) // Add a delay if needed for the producer to finish

		ginkgo.By("Executing queue stat command in Artemis pod")
		output := execCommandInPod(kubeClient, namespace, podName,
			"./amq-broker/bin/artemis queue stat --url tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616 --user cgi --password cgi --maxRows 200 --clustered")

		ginkgo.By("Checking if 'TEST' queue has 300 in its row")
		lines := strings.Split(output, "\n")
		found := false
		for _, line := range lines {
			if strings.Contains(line, "TEST") && strings.Contains(line, "300") {
				found = true
				break
			}
		}

		gomega.Expect(found).To(gomega.BeTrue(), fmt.Sprintf("Expected 'TEST' queue with 300 not found in output:\n%s", output))
	})
})

func execCommandInPod(clientset *kubernetes.Clientset, namespace, podName, command string) string {
	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		Param("container", "ex-aao-ss-0").
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "true").
		Param("command", "/bin/bash").
		Param("-c", command)

	executor := &exec.PodExecutor{}
	executor.StreamOptions.Tty = true
	executor.StreamOptions.IOStreams.Out = ginkgo.GinkgoWriter
	executor.StreamOptions.IOStreams.ErrOut = ginkgo.GinkgoWriter

	err := executor.StreamOptions.Validate()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	streamOptions := executor.StreamOptions.Copy()
	transport, upgrader, err := spdy.RoundTripperFor(streamOptions)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	request := req.URL()
	request = request.Scheme("https")

	streamOptions.Upgrade = upgrader.Upgrade
	streamOptions.Request = request

	podStreamOptions := &exec.PodStreamOptions{
		IOStreams: streamOptions.IOStreams,
		Tty:       streamOptions.Tty,
	}

	streamOptions.Upgrade.Dial = remote.NewDialFunc(transport.Dial)

	podStreamOptions.StreamOptions = streamOptions

	fn := executor.StreamOptions.Upgrade.RoundTripper

	_, _, err = remote.StreamWithOptions(podStreamOptions, fn)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	return "" // Modify the return type based on your use case
}

func TestArtemis(t *testing.T) {
	ginkgo.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Artemis Suite")
}
