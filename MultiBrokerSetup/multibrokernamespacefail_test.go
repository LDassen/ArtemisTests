package MultiBrokerSetup_test

import (
	"path/filepath"

	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ = BeforeSuite(func() {
	// Set up Kubernetes client using in-cluster configuration
	config, err := rest.InClusterConfig()
	Expect(err).NotTo(HaveOccurred())

	// You can use an underscore (_) to indicate that the variable is intentionally unused
	_, err = kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Deploying to Non-existing Namespace", func() {
	It("Should fail to deploy in a non-existing namespace", func() {
		namespace := "nonexistent-namespace"

		// Get the current directory
		currentDir, err := filepath.Abs(filepath.Dir("."))
		Expect(err).NotTo(HaveOccurred())

		// Construct the full path to the deployment file
		deploymentFile := filepath.Join(currentDir, "ex-aao.yaml")

		// Use kubectl-testkube library to apply manifests
		err = commands.Apply(deploymentFile, "--namespace="+namespace)
		Expect(err).To(HaveOccurred()) // Expect an error as the namespace is non-existing
	})
})

var _ = AfterSuite(func() {
	// Clean up resources if needed
})
