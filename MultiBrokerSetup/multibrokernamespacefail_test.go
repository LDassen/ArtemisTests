package MultiBrokerSetup_test

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

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

	// Download kubectl
	cmd := exec.Command("curl", "-LO", "https://storage.googleapis.com/kubernetes-release/release/v1.23.0/bin/linux/amd64/kubectl")
	err = cmd.Run()
	Expect(err).NotTo(HaveOccurred())

	// Move kubectl to the desired location
	kubectlPath := "/usr/local/bin/kubectl"
	_, err = exec.Command("mv", "kubectl", kubectlPath).CombinedOutput()
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

		// Use the downloaded kubectl binary within the pod
		kubectlPath := "/path/to/kubectl"
		cmd := exec.Command(kubectlPath, "apply", "-f", deploymentFile, "--namespace="+namespace)
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

// downloadKubectl downloads kubectl binary if not already present
func downloadKubectl() error {
	// Check if kubectl is already present
	kubectlPath, err := exec.LookPath("kubectl")
	if err == nil {
		fmt.Println("kubectl is already present:", kubectlPath)
		return nil
	}

	// Download kubectl based on the OS
	var downloadURL string
	switch runtime.GOOS {
	case "linux":
		downloadURL = "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
	case "darwin":
		downloadURL = "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/darwin/amd64/kubectl"
	case "windows":
		downloadURL = "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/windows/amd64/kubectl.exe"
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Download kubectl
	downloadCmd := exec.Command("sh", "-c", fmt.Sprintf("curl -LO %s && chmod +x kubectl && mv kubectl /path/to/", downloadURL))
	downloadOutput, downloadErr := downloadCmd.CombinedOutput()
	if downloadErr != nil {
		return fmt.Errorf("failed to download kubectl: %v\nOutput:\n%s", downloadErr, string(downloadOutput))
	}

	return nil
}
