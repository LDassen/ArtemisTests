package MultiBrokerSetup_test

import (
	"fmt"
	"os"
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

		// Get the current directory
		currentDir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())

		// Construct the full path to the deployment file
		deploymentFile := currentDir + "/ex-aao.yaml"

		// Apply the deployment file in the non-existing namespace
		cmd := exec.Command("kubectl", "apply", "-f", deploymentFile, "--namespace="+namespace)
		output, err := cmd.CombinedOutput()

		// Print both stdout and stderr for debugging
		fmt.Println("Command stdout:", string(output))
		fmt.Println("Command stderr:", err)

		// Verify that the error indicates a non-existing namespace
		Expect(err).To(HaveOccurred())
		Expect(string(output)).To(ContainSubstring("does not match the namespace \"" + namespace + "\""))
	})
})

var _ = AfterSuite(func() {
	// Clean up resources if needed
})
