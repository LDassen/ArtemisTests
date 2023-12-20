package MultiBrokerSetup

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ = ginkgo.Describe("Artemis Broker Deployment", func() {
	var kubeClient *kubernetes.Clientset

	ginkgo.BeforeEach(func() {
		// Load the in-cluster Kubernetes config
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Create a Kubernetes client
		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
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

func TestArtemisBrokerDeployment(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Artemis Broker Deployment Suite")
}
