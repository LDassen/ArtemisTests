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

	ginkgo.It("Should fail to apply a deployment file for Artemis to a non-existing namespace", func() {
		fileName := "case_2-ex-aao.yaml"
		namespace := "non-existing"

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
		gomega.Expect(err).To(gomega.MatchError("namespaces \"non-existing\" not found"))
	})

	ginkgo.AfterEach(func() {
		// Cleanup logic if needed
	})
})
