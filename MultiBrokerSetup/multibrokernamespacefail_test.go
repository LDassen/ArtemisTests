package artemistests

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os/exec"
	"testing"
)

var _ = BeforeSuite(func() {
	// Set up any prerequisites before the test suite
})

var _ = AfterSuite(func() {
	// Clean up after the test suite
})

var _ = Describe("Artemis Deployment", func() {
	It("should fail to apply deployment to a non-existing namespace", func() {
		// Your test logic here

		// Example: Get the path to the deployment YAML file in the same directory
		deploymentFile := filepath.Join(".", "deployment.yaml")
		namespace := "non-existing-namespace"
		cmd := exec.Command("kubectl", "apply", "-f", deploymentFile, "--namespace="+namespace)
		err := cmd.Run()

		// Expect an error to occur
		Expect(err).To(HaveOccurred(), "Expected an error but got none")
	})

	// Add more test cases as needed
})

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}
