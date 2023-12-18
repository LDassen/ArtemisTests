package AutoCreationQueue_test

import (
	"bytes"
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Artemis Broker", func() {

	It("should run a command inside the Artemis broker", func() {
		// Create Kubernetes client
		clientset, err := createKubernetesClient()
		Expect(err).NotTo(HaveOccurred())

		// Replace this command with the actual command you want to run inside the broker
		commandToRun := "./amq-broker/bin/artemis producer --user cgi --password cgi --url tcp://10.204.0.39:61616 --message-count 100"

		// Set the namespace and pod name directly
		namespace := "activemq-artemis-brokers"
		podName := "ex-aao-ss-0"

		// Run the command inside the specific Artemis broker pod
		output, err := runCommandInsideKubernetesPod(clientset, podName, namespace, commandToRun)
		Expect(err).NotTo(HaveOccurred())

		// Add your assertions based on the command output
		Expect(output).To(ContainSubstring("expected-output"))
	})
})

// Helper function to create a Kubernetes client
func createKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// Helper function to run a command inside a Kubernetes pod using exec
func runCommandInsideKubernetesPod(clientset *kubernetes.Clientset, podName, namespace, command string) (string, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	// Create an exec command
	execCommand := exec.Command("/bin/bash", "-c", command)

	// Set the correct working directory
	execCommand.Dir = "/home/jboss"

	// Capture the command output
	var stdout, stderr bytes.Buffer
	execCommand.Stdout = &stdout
	execCommand.Stderr = &stderr

	// Run the exec command
	err = execCommand.Run()
	if err != nil {
		return "", fmt.Errorf("error executing command: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	return stdout.String(), nil
}

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}
