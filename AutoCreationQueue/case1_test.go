package AutoCreationQueue_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/wait"
	"k8s.io/client-go/util/wait/remote"
	"k8s.io/client-go/util/wait/transport"
	"k8s.io/client-go/util/wait/transport/spdy"
	"k8s.io/client-go/util/wait/transport/ws"
	"k8s.io/client-go/util/yaml"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/exec"
	"k8s.io/client-go/util/wait"
	"k8s.io/kubernetes/pkg/kubectl/exec"
)

var _ = Describe("Artemis Test", func() {
	var (
		kubeClient *kubernetes.Clientset
		namespace  = "activemq-artemis-brokers"
		ctx        = context.TODO()
	)

	BeforeSuite(func() {
		// Initialize Kubernetes client
		config, err := rest.InClusterConfig()
		if err != nil {
			home := homedir.HomeDir()
			kubeconfig := filepath.Join(home, ".kube", "config")
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			Expect(err).NotTo(HaveOccurred())
		}
		kubeClient, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should test Artemis functionality", func() {
		By("Searching for Artemis pods in the namespace")
		podList, err := kubeClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(podList.Items).NotTo(BeEmpty())

		podName := "ex-aao-ss-0"

		By("Executing command in Artemis pod")
		execCommandInPod(kubeClient, namespace, podName,
			"./amq-broker/bin/artemis producer --user cgi --password cgi --url tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616 --message-count 100 --user cgi --password cgi --maxRows 200")

		time.Sleep(5 * time.Second) // Add a delay if needed for the producer to finish

		By("Executing queue stat command in Artemis pod")
		output := execCommandInPod(kubeClient, namespace, podName,
			"./amq-broker/bin/artemis queue stat --url tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616 --user cgi --password cgi --maxRows 200 --clustered")

		By("Checking if 'TEST' queue has 300 in its row")
		lines := strings.Split(output, "\n")
		found := false
		for _, line := range lines {
			if strings.Contains(line, "TEST") && strings.Contains(line, "300") {
				found = true
				break
			}
		}

		Expect(found).To(BeTrue(), fmt.Sprintf("Expected 'TEST' queue with 300 not found in output:\n%s", output))
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
	executor.StreamOptions.IOStreams.Out = GinkgoWriter
	executor.StreamOptions.IOStreams.ErrOut = GinkgoWriter

	err := executor.StreamOptions.Validate()
	Expect(err).NotTo(HaveOccurred())

	streamOptions := executor.StreamOptions.Copy()
	transport, upgrader, err := spdy.RoundTripperFor(streamOptions)
	Expect(err).NotTo(HaveOccurred())

	request, err := req.URL()
	Expect(err).NotTo(HaveOccurred())
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

	return wait.Poll(100*time.Millisecond, wait.ForeverTestTimeout, func() (bool, error) {
		_, _, err := remote.StreamWithOptions(podStreamOptions, fn)
		if err == transport.ErrStreamPrefixNotFound {
			return false, nil
		}
		return true, err
	}).Until(func() (string, error) {
		return "", nil
	})
}

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}
