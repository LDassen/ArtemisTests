package test_test

import (
	"os"
	"testing"
	"path/filepath"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = ginkgo.Describe("Artemis Broker Setup", func() {
	ginkgo.It("should have three brokers running", func() {
		config, err := loadKubeConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		clientset, err := kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		podList, err := clientset.CoreV1().Pods("default").List(ginkgo.Context, metav1.ListOptions{
			LabelSelector: "app=artemis-broker",
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		gomega.Expect(len(podList.Items)).To(gomega.Equal(3), "Expected 3 Artemis brokers, but found %d", len(podList.Items))
	})
})

func TestArtemis(t *testing.T) {
	ginkgo.RunSpecs(t, "Artemis Suite")
}

func loadKubeConfig() (*rest.Config, error) {
	var kubeconfig string
	if home := homeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return nil, fmt.Errorf("home directory not found")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
