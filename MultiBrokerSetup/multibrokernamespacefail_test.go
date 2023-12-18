package MultiBrokerSetup_test

import (
	"context"
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/api/errors" // Import the correct package
)

var _ = Describe("Deploying to Non-existing Namespace", func() {
	var clientset *kubernetes.Clientset

	BeforeSuite(func() {
		// Set up Kubernetes client using in-cluster configuration
		config, err := rest.InClusterConfig()
		Expect(err).NotTo(HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("Should fail to deploy in a non-existing namespace", func() {
		namespace := "nonexistent-namespace"
		deploymentFile := "ex-aao.yaml"

		// Apply the deployment file in the non-existing namespace
		cmd := exec.Command("kubectl", "apply", "-f", deploymentFile, "--namespace="+namespace)
		output, err := cmd.CombinedOutput()

		// Verify that the error indicates a non-existing namespace
		Expect(err).To(HaveOccurred())
		Expect(output).To(ContainSubstring(fmt.Sprintf("namespace %s not found", namespace)))

		// Alternatively, you can use the Kubernetes client to check if the namespace exists
		_, err = clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
		Expect(err).To(HaveOccurred())
		Expect(kubernetes.IsNotFound(err)).To(BeTrue())
	})

	AfterSuite(func() {
		// Clean up resources if needed
	})
})
