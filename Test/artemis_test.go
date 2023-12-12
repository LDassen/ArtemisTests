package test_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/rest"
	"k8s.io/client-go/util/wait"
)

var _ = Describe("Artemis Broker Pods", func() {
	It("should have the correct number of 'broker' pods running", func() {
		namespace := "activemq-artemis-operator"
		expectedPodCount := 3 // Set your expected number of 'broker' pods

		// Get Kubernetes client
		client, err := getClient()
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		// Get the list of pods in the specified namespace with a label selector
		pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "application=ex-aao-app",
		})
		Expect(err).To(BeNil(), "Error getting pods: %v", err)

		// Count the number of pods
		actualPodCount := len(pods.Items)

		// Assert that the actual count matches the expected count
		Expect(actualPodCount).To(Equal(expectedPodCount), "Expected %d 'broker' pods, but found %d", expectedPodCount, actualPodCount)
	})
})

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}

func getClient() (*kubernetes.Clientset, error) {
	// Use in-cluster config if available, else use the default kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	// Wait for the Kubernetes client to be ready
	err = wait.PollImmediate(5*time.Second, 1*time.Minute, func() (bool, error) {
		return clientset.Discovery().ServerPreferredResources()
	})

	return clientset, err
}

