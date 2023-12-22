package MultiBrokerSetup_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Kubernetes Apply Deployment Test", func() {
	var clientset *kubernetes.Clientset

	ginkgo.BeforeEach(func() {
		// Set up the Kubernetes client
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Should NOT fail to apply a deployment file for Artemis to an existing namespace", func() {
		// Filename of the deployment YAML file
		fileName := "ex-aao.yaml"

		// Namespace where the deployment will be applied
		namespace := "activemq-artemis-brokers"

		// Read the deployment file
		filePath, err := filepath.Abs(fileName)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		fileBytes, err := ioutil.ReadFile(filePath)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Decode the YAML manifest
		decode := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(fileBytes), 1024)
		var deployment appsv1.Deployment
		err = decode.Decode(&deployment)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Apply the deployment to the existing namespace
		_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error creating deployment: %v", err)
	})

	ginkgo.AfterEach(func() {
        namespace := "activemq-artemis-brokers"
        deploymentName := "ex-aao" // Ensure this matches the deployment name in your YAML file

        // Delete the deployment
        err := clientset.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
        if err != nil {
            fmt.Printf("Error deleting deployment: %v\n", err)
        }

        // Delete the broker pods
        deleteOptions := metav1.DeleteOptions{}
        labelSelector := "application=ex-aao-app"
        err = clientset.CoreV1().Pods(namespace).DeleteCollection(context.TODO(), deleteOptions, metav1.ListOptions{
            LabelSelector: labelSelector,
        })
        if err != nil {
            fmt.Printf("Error deleting broker pods: %v\n", err)
        }
    })

	ginkgo.It("should have 3 'broker' pods running in the namespace with the app label 'application=ex-aao-app'", func() {
		// Namespace where the pods are expected to be running
		namespace := "activemq-artemis-brokers"

		// Expected number of pods
		expectedPodCount := 3

		// Fetch the pods with the specified label in the namespace
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "application=ex-aao-app",
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error getting pods: %v", err)

		// Debugging statements
		fmt.Printf("Retrieved %d pods in namespace %s with the label 'application=ex-aao-app'\n", len(pods.Items), namespace)
		for _, pod := range pods.Items {
			fmt.Printf("Pod Name: %s\n", pod.Name)
			// Add more details as needed
		}

		// Actual number of pods
		actualPodCount := len(pods.Items)

		// Assertion: Check if the actual number matches the expected number
		gomega.Expect(actualPodCount).To(gomega.Equal(expectedPodCount),
			"Expected %d 'broker' pods with the label 'application=ex-aao-app', but found %d", expectedPodCount, actualPodCount)
	})
})
