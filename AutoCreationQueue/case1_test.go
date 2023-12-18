package AutoCreationQueue_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
			kubeConfigPath := os.Getenv("KUBECONFIG")
			config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
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

		podName := podList.Items[0].Name

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
		Param("command", "/bin/bash", "-c", command)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	Expect(err).NotTo(HaveOccurred())

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		IOStreams: remotecommand.IOStreams{
			Out: &stdout,
			Err: &stderr,
		},
		Tty: true,
	})
	Expect(err).NotTo(HaveOccurred())

	return stdout.String()
}

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}
