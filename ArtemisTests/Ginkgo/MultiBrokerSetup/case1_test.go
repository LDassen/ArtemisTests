package MultiBrokerSetup_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"

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
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		clientset, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("should apply a deployment file for Artemis to a namespace", func() {
		fileName := "case_1-ex-aao.yaml"
		namespace := "activemq-artemis-brokers" // Replace with your existing namespace or a new one

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

		// Create or get the namespace
		_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}, metav1.CreateOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Apply the deployment to the namespace
		_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), &deployment, metav1.CreateOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.AfterEach(func() {
		// Cleanup logic if needed
	})
})
