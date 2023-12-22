package MultiBrokerSetup_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/util/wait"
)

var _ = ginkgo.Describe("Kubernetes Apply Deployment Test", func() {
	var clientset *kubernetes.Clientset

	ginkgo.BeforeEach(func() {
		// Set up the Kubernetes client
		config, err := rest.InClusterConfig()
		gomega.Expect(err).To(gomega.BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).To(gomega.BeNil(), "Error creating Kubernetes client: %v", err)
		gomega.RegisterFailHandler(ginkgo.Fail)
	})

	ginkgo.AfterEach(func() {
		// Cleanup: Delete the deployment
		err := clientset.AppsV1().Deployments("activemq-artemis-brokers").Delete(context.TODO(), "ex-aao", metav1.DeleteOptions{})
		gomega.Expect(err).To(gomega.BeNil(), "Error deleting deployment: %v", err)
	})

	ginkgo.It("should apply a deployment file for Artemis to a namespace and have the correct number of 'broker' pods running", func() {
		fileName := "ex-aao.yaml"
		namespace := "activemq-artemis-brokers"

		// Apply the deployment to the namespace
		err := applyDeploymentFromFile(clientset, fileName, namespace)
		gomega.Expect(err).To(gomega.BeNil(), "Error applying deployment: %v", err)

		// Wait for the deployment to be available
		ginkgo.By("Waiting for the deployment to be available")
		err = waitForDeployment(clientset, namespace, "ex-aao", 3, 5*time.Minute) // Adjust timeout as needed
		gomega.Expect(err).To(gomega.BeNil(), "Error waiting for deployment: %v", err)

		// Check the number of 'broker' pods
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "app=broker",
		})
		gomega.Expect(err).To(gomega.BeNil(), "Error getting pods: %v", err)

		// Debugging statements
		fmt.Printf("Retrieved %d pods in namespace %s\n", len(pods.Items), namespace)
		for _, pod := range pods.Items {
			fmt.Printf("Pod Name: %s\n", pod.Name)
			// Add more details as needed
		}

		expectedPodCount := 3 // Set your expected number of 'broker' pods
		actualPodCount := len(pods.Items)

		gomega.Expect(actualPodCount).To(gomega.Equal(expectedPodCount), "Expected %d 'broker' pods, but found %d", expectedPodCount, actualPodCount)
	})
})

// Helper function to apply a deployment from a file
func applyDeploymentFromFile(clientset *kubernetes.Clientset, fileName, namespace string) error {
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Decode the YAML manifest
	decode := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(fileBytes), 1024)
	var deployment appsv1.Deployment
	err = decode.Decode(&deployment)
	if err != nil {
		return err
	}

	// Apply the deployment to the namespace
	_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
	return err
}