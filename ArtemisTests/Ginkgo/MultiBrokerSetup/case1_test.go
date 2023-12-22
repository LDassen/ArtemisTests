package MultiBrokerSetup_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	ginkgo.It("should apply a deployment file for Artemis to a namespace", func() {
		fileName := "ex-aao.yaml"
		namespace := "activemq-artemis-brokers"

		// Read the file
		filePath, err := filepath.Abs(fileName)
		gomega.Expect(err).To(gomega.BeNil(), "Error getting absolute file path: %v", err)

		fileBytes, err := ioutil.ReadFile(filePath)
		gomega.Expect(err).To(gomega.BeNil(), "Error reading file: %v", err)

		// Decode the YAML manifest
		decode := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(fileBytes), 1024)
		var deployment appsv1.Deployment
		err = decode.Decode(&deployment)
		gomega.Expect(err).To(gomega.BeNil(), "Error decoding YAML: %v", err)

		// Apply the deployment to the namespace
		_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.BeNil(), "Error creating deployment: %v", err)

		// Wait for the deployment to be available
		ginkgo.By("Waiting for the deployment to be available")
		err = waitForDeployment(clientset, namespace, deployment.Name, 3, 5*time.Minute) // Adjust timeout as needed
		gomega.Expect(err).To(gomega.BeNil(), "Error waiting for deployment: %v", err)
	})

	ginkgo.AfterEach(func() {
		// Cleanup logic if needed
	})

	ginkgo.It("should have the correct number of 'broker' pods running", func() {
		config, err := rest.InClusterConfig()
		gomega.Expect(err).To(gomega.BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		gomega.Expect(err).To(gomega.BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artemis-brokers"
		expectedPodCount := 3 // Set your expected number of 'broker' pods

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: "application=ex-aao-app"})
		gomega.Expect(err).To(gomega.BeNil(), "Error getting pods: %v", err)

		// Debugging statements
		fmt.Printf("Retrieved %d pods in namespace %s\n", len(pods.Items), namespace)
		for _, pod := range pods.Items {
			fmt.Printf("Pod Name: %s\n", pod.Name)
			// Add more details as needed
		}

		actualPodCount := len(pods.Items)

		gomega.Expect(actualPodCount).To(gomega.Equal(expectedPodCount), "Expected %d 'broker' pods, but found %d", expectedPodCount, actualPodCount)
	})
})

// Helper function to wait for the deployment to be available
func waitForDeployment(clientset *kubernetes.Clientset, namespace, deploymentName string, replicas int32, timeout time.Duration) error {
	return wait.PollImmediate(10*time.Second, timeout, func() (bool, error) {
		deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if deployment.Status.AvailableReplicas == replicas {
			return true, nil
		}

		return false, nil
	})
}
