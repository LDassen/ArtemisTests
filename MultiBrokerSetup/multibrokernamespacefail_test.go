package MultiBrokerSetup_test

import (
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apply Kubernetes Configuration File and Get Error Logs", func() {
	It("should apply a configuration file and retrieve error logs", func() {
		namespace := "unexisting-namespace"
		configFilePath := "ex-aao.yaml" // Change this to the path of your single configuration file

		// Execute kubectl apply command
		cmd := exec.Command("kubectl", "apply", "-f", configFilePath, "--namespace", namespace)
		output, err := cmd.CombinedOutput()

		// Check for errors
		Expect(err).ToNot(BeNil(), "Expected an error applying the configuration file")

		// Print error and kubectl output
		fmt.Printf("Error applying the configuration file: %v\n", err)
		fmt.Printf("kubectl output:\n%s\n", string(output))
	})
})
