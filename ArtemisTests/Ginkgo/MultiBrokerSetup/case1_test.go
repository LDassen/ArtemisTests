package MultiBrokerSetup_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"time"
	
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	ginkgo.It("should apply a deployment file for Artemis to a namespace", func() {
		fileName := "ex-aao.yaml" //"case_1-ex-aao.yaml"
		namespace := "activemq-artemis-brokers" // Replace with your existing namespace or a new one

		// Check if the namespace already exists
		_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
		if err != nil && errors.IsNotFound(err) {
			// The namespace does not exist, so create it
			_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}, metav1.CreateOptions{})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		} else if err != nil {
			// Handle other errors, if any
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		}

		// Read the file
		filePath, err := filepath.Abs(fileName)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		fileBytes, err := ioutil.ReadFile(filePath)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Decode the YAML manifest
		decode := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(fileBytes), 1024)
		var deployment appsv1.Deployment
		err = decode.Decode(&deployment)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Apply the deployment to the namespace
		_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for the deployment to be available
		ginkgo.By("Waiting for the deployment to be available")
		err = waitForDeployment(clientset, namespace, deployment.Name, 3, 5*time.Minute) // Adjust timeout as needed
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterEach(func() {
		// Cleanup logic if needed
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