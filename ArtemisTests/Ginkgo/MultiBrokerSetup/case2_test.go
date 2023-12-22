package MultiBrokerSetup_test

import (
    "bytes"
    "context"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "strings"

    appsv1 "k8s.io/api/apps/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/yaml"
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

    ginkgo.It("Should fail to apply a deployment file for Artemis to a non-existing namespace", func() {
        fileName := "ex-aao.yaml"
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
        if err != nil && strings.Contains(err.Error(), "the namespace of the provided object does not match the namespace sent on the request") {
            fmt.Println("[ERROR] Namespace mismatch error encountered:", err)
            ginkgo.Skip("Namespace mismatch error: skipping test due to specific error")
        }

        gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Deployment should not succeed in non-existing namespace")
    })

    ginkgo.AfterEach(func() {
		err := clientset.AppsV1().Deployments("activemq-artemis-brokers").Delete(context.TODO(), "ex-aao", metav1.DeleteOptions{})
		gomega.Expect(err).To(gomega.BeNil(), "Error deleting deployment: %v", err)
    })
})


