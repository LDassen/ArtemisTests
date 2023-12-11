package test_test

import (
	"fmt"
	"os"
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"test"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("Pod Count in Namespace", func() {
	var (
		clientset *kubernetes.Clientset
		namespace string
	)

	BeforeSuite(func() {
		// Load Kubernetes config from default location or provide your kubeconfig path
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.Getenv("HOME") + "/.kube/config"
		}

		// Build Kubernetes client
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		Expect(err).NotTo(HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		// Set the namespace you want to test
		namespace = "your-namespace"
	})

	It("should have the expected number of pods", func() {
		podList, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		// Define the expected number of pods
		expectedPodCount := 3

		// Assert the actual number of pods matches the expected count
		Expect(len(podList.Items)).To(Equal(expectedPodCount))
	})
})

func TestPodCount(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pod Count Suite")
}