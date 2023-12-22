package MultiBrokerSetup_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
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

	ginkgo.It("Should NOT fail to apply a deployment file for Artemis to a non-existing namespace", func() {
		fileName := "ex-aao.yaml"
		namespace := "activemq-artemis-brokers"

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

		// Try to apply the deployment to the non-existing namespace
		_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
	})

	ginkgo.AfterEach(func() {
		err := clientset.AppsV1().Deployments("activemq-artemis-brokers").Delete(context.TODO(), "ex-aao", metav1.DeleteOptions{})
		gomega.Expect(err).To(gomega.BeNil(), "Error deleting deployment: %v", err)
	})

	ginkgo.It("should have 3 'broker' pods running in the namespace with the app label 'application=ex-aao-app'", func() {
		namespace := "activemq-artemis-brokers"
		expectedPodCount := 3

		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "application=ex-aao-app",
		})
		gomega.Expect(err).To(gomega.BeNil(), "Error getting pods: %v", err)

		// Debugging statements
		fmt.Printf("Retrieved %d pods in namespace %s with the label 'application=ex-aao-app'\n", len(pods.Items), namespace)
		for _, pod := range pods.Items {
			fmt.Printf("Pod Name: %s\n", pod.Name)
			// Add more details as needed
		}

		actualPodCount := len(pods.Items)

		gomega.Expect(actualPodCount).To(gomega.Equal(expectedPodCount),
			"Expected %d 'broker' pods with the label 'application=ex-aao-app', but found %d", expectedPodCount, actualPodCount)
	})
})
