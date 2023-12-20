package MultiBrokerSetup

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Artemis Broker Deployment", func() {
	var kubeConfigPath string

	ginkgo.BeforeEach(func() {
		// Find kubeconfig file path
		kubeConfigPath = findKubeConfig()
		gomega.Expect(kubeConfigPath).NotTo(gomega.BeEmpty())
	})

	ginkgo.It("should apply Artemis broker deployment file to a non-existing namespace", func() {
		// Set the path to your Artemis broker deployment file
		artemisDeploymentFilePath := "ex-aao.yaml"

		// Set the non-existing namespace
		nonExistingNamespace := "non-existing-namespace"

		// Build the kubectl apply command
		cmd := exec.Command("kubectl", "apply", "-f", artemisDeploymentFilePath, "--namespace="+nonExistingNamespace)

		// Run the command
		output, err := cmd.CombinedOutput()
		gomega.Expect(err).To(gomega.HaveOccurred())

		// Check if the error message contains the expected information
		expectedErrorMessage := "namespace \"" + nonExistingNamespace + "\" not found"
		gomega.Expect(strings.Contains(string(output), expectedErrorMessage)).To(gomega.BeTrue())
	})

	ginkgo.AfterEach(func() {
		// Cleanup if necessary
	})

})

func findKubeConfig() string {
	// Check KUBECONFIG environment variable first
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig != "" {
		return kubeconfig
	}

	// If KUBECONFIG is not set, use the default kubeconfig file location
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".kube", "config")
}

func TestArtemisBrokerDeployment(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Artemis Broker Deployment Suite")
}