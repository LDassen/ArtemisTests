package MultiBrokerSetup

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/kubeconfig"
	"os"
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
		deploymentFile := filepath.Join(".", "ex-aao.yaml")
		namespace := "nonexistent-namespace"

		// Load kubeconfig
		home := homedir.HomeDir()
		config, err := kubeconfig.LoadKubeConfig(filepath.Join(home, ".kube", "config"))
		Expect(err).NotTo(HaveOccurred(), "Error loading kubeconfig")

		// Create Kubernetes client
		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred(), "Error creating Kubernetes client")

		// Read deployment YAML file
		deployment, err := clientset.AppsV1().Deployments(namespace).CreateFromFile(deploymentFile)
		Expect(err).To(HaveOccurred(), "Expected an error but got none")

		// Check for the specific error message in the error
		expectedErrorMessage := "the namespace from the provided object"
		Expect(strings.Contains(err.Error(), expectedErrorMessage)).To(BeTrue(),
			"Expected error message not found in error: "+err.Error())

		// Optionally, you can assert other conditions based on your needs
		Expect(deployment).To(BeNil(), "Deployment should be nil due to error")
	})

	// Add more test cases as needed
})

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}
