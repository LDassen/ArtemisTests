package AutoCreationQueue_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Artemis Broker", func() {

	It("should run a command inside the Artemis broker", func() {
		// Replace these values with your actual Kubernetes configuration
		podName := "artemis-broker-pod-name"
		namespace := "default"

		// Replace this command with the actual command you want to run inside the broker
		commandToRun := "./amq-broker/bin/artemis producer --user cgi --password cgi --url tcp://10.204.0.39:61616 --message-count 100"

		// Run the command inside the specific Artemis broker pod
		output, err := runCommandInsideKubernetesPod(podName, namespace, commandToRun)
		Expect(err).NotTo(HaveOccurred())

		// Add your assertions based on the command output
		Expect(output).To(ContainSubstring("expected-output"))
	})
})

// Helper function to run a command inside a Kubernetes pod using exec
func runCommandInsideKubernetesPod(podName, namespace, command string) (string, error) {
	// Create an exec command
	execCommand := exec.Command("kubectl", "exec", "-it", "pod", podName, namespace, "--", "/bin/bash", command)

	// Capture the command output
	var stdout, stderr bytes.Buffer
	execCommand.Stdout = &stdout
	execCommand.Stderr = &stderr

	// Run the exec command
	err := execCommand.Run()
	if err != nil {
		return "", fmt.Errorf("error executing command: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	return stdout.String(), nil
}

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}
