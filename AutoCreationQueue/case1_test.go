package AutoCreationQueue_test

import (
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Artemis Broker", func() {
	var (
		brokerCmd *exec.Cmd
	)

	BeforeSuite(func() {
		// Start the Artemis broker as a separate process
		brokerCmd = exec.Command("artemis", "run")
		err := brokerCmd.Start()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterSuite(func() {
		// Stop the Artemis broker after all tests are finished
		err := brokerCmd.Process.Kill()
		Expect(err).NotTo(HaveOccurred())
	})

	It("should run a command inside the Artemis broker", func() {
		// Replace this command with the actual command you want to run inside the broker
		commandToRun := "your-command-here"

		// Run the command inside the Artemis broker
		output, err := runCommandInsideBroker(commandToRun)
		Expect(err).NotTo(HaveOccurred())

		// Add your assertions based on the command output
		Expect(output).To(ContainSubstring("expected-output"))
	})
})

// Helper function to run a command inside the Artemis broker
func runCommandInsideBroker(command string) (string, error) {
	// Replace this with the actual Artemis broker URL and credentials
	brokerURL := "tcp://localhost:61616"
	username := "your-username"
	password := "your-password"

	// Construct the Artemis command to run a command inside the broker
	artemisCmd := exec.Command("artemis", "producer", "--message", command, "--destination", "exampleQueue", "--url", brokerURL, "--user", username, "--password", password)

	// Run the Artemis command and capture the output
	output, err := artemisCmd.CombinedOutput()
	return string(output), err
}

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}

