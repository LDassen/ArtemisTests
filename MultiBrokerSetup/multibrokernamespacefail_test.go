package MultiBrokerSetup_test

import (
	"os/exec"

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
		deploymentFile := "ex-aao.yaml"

		// Apply the deployment file in the non-existing namespace
		cmd := exec.Command("kubectl", "apply", "-f", deploymentFile, "--namespace="+namespace)
		output, err := cmd.CombinedOutput()

		// Verify that the error indicates a non-existing namespace
		Expect(err).To(HaveOccurred())
		Expect(output).To(ContainSubstring("does not match the namespace \"" + namespace + "\""))

		// Print the command output for debugging
		fmt.Println("Command output:", string(output))
	})
})

var _ = AfterSuite(func() {
	// Clean up resources if needed
})
