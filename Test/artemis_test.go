package test_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var _ = ginkgo.Describe("Artemis Broker Setup", func() {
	g := gomega.NewGomegaWithT(ginkgo.GinkgoT())

	// Your test goes here
	ginkgo.It("should have three brokers running", func() {
		// Load Kubernetes config
		config, err := loadKubeConfig()
		g.Expect(err).NotTo(gomega.HaveOccurred())

		// Create Kubernetes client
		clientset, err := kubernetes.NewForConfig(config)
		g.Expect(err).NotTo(gomega.HaveOccurred())

		// Get Artemis broker pods
		podList, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
			LabelSelector: "app=artemis-broker",
		})
		g.Expect(err).NotTo(gomega.HaveOccurred())

		// Assert that there are exactly 3 Artemis broker pods
		g.Expect(len(podList.Items)).To(gomega.Equal(3), "Expected 3 Artemis brokers, but found %d", len(podList.Items))
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

	config, err := rest.InClusterConfig()
	if err != nil {
		// If not running inside a cluster, use kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
