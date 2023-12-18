package AutoCreationQueue_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Artemis Broker", func() {

	It("should run a command inside the Artemis broker", func() {
		// Replace this command with the actual command you want to run inside the broker
		commandToRun := "./amq-broker/bin/artemis producer --user cgi --password cgi --url tcp://7.65.87.2:61616 --message-count 100"

		// Run the command inside the specific Artemis broker pod using kubectl exec
		output, err := runCommandInsideKubernetesPod("artemis-broker-pod-name", commandToRun)
		Expect(err).NotTo(HaveOccurred())

		// Add your assertions based on the command output
		Expect(output).To(ContainSubstring("expected-output"))
	})
})

// Helper function to run a command inside a Kubernetes pod
func runCommandInsideKubernetesPod(podName, command string) (string, error) {
	// Construct the kubectl exec command to run a command inside the pod
	kubectlCmd := exec.Command("kubectl", "exec", "-it", "pod/"podName, "--", "/bin/bash", command)

	// Run the kubectl exec command and capture the output
	output, err := kubectlCmd.CombinedOutput()
	return string(output), err
}

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}
