package MessageMigration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("ActiveMQ Artemis Message Migration Test", func() {
	var dynamicClient dynamic.Interface
	var kubeClient *kubernetes.Clientset
	var namespace string
	var resourceGVR schema.GroupVersionResource

	ginkgo.BeforeEach(func() {
		var err error
		var kubeconfig *rest.Config

		// Use in-cluster config if running in a Kubernetes cluster
		if kubeconfig, err = rest.InClusterConfig(); err != nil {
			// If not in a cluster, use kubeconfig file from home directory
			home := homedir.HomeDir()
			kubeconfig, err = clientcmd.BuildConfigFromFlags("", filepath.Join(home, ".kube", "config"))
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		}

		// Create dynamic client
		dynamicClient, err = dynamic.NewForConfig(kubeconfig)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Create Kubernetes client
		kubeClient, err = kubernetes.NewForConfig(kubeconfig)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		namespace = "activemq-artemis-brokers"
		resourceGVR = schema.GroupVersionResource{
			Group:    "broker.amq.io",
			Version:  "v1beta1",
			Resource: "activemqartemises",
		}
	})

	ginkgo.It("Should scale down ActiveMQ Artemis brokers and verify message migration", func() {
		// Scale down ActiveMQ Artemis brokers
		scaleDownSize := 2
		scaleDownActiveMQArtemis(scaleDownSize)

		// Check if brokers scaled down correctly
		verifyBrokerScaling(scaleDownSize)

		// Wait for the drainer pod to start and print its logs
		drainerPodLogs := waitForDrainerPod()
		fmt.Printf("Drainer Pod Logs:\n%s\n", drainerPodLogs)

		// TODO: Add logic to verify message migration in the logs
		// Example: assertMessageMigration(drainerPodLogs)
	})
})

// Helper function to scale down ActiveMQ Artemis brokers
func scaleDownActiveMQArtemis(size int) {
	fileName := "ex-aaoMM.yaml"

	filePath, err := filepath.Abs(fileName)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	fileBytes, err := ioutil.ReadFile(filePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, _, err = decUnstructured.Decode(fileBytes, nil, obj)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	obj.SetAPIVersion("broker.amq.io/v1beta1")
	obj.SetKind("ActiveMQArtemis")

	resourceClient := dynamicClient.Resource(resourceGVR).Namespace(namespace)

	// Set replicas to the desired size
	unstructured.SetNestedField(obj.Object, int64(size), "spec", "replicas")

	_, err = resourceClient.Update(context.TODO(), obj, metav1.UpdateOptions{})
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error scaling down ActiveMQ Artemis brokers")
}

// Helper function to verify if brokers scaled down correctly
func verifyBrokerScaling(expectedSize int) {
	pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=activemq-artemis-broker",
	})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	actualSize := len(pods.Items)
	gomega.Expect(actualSize).To(gomega.Equal(expectedSize), "Brokers did not scale down correctly")
}

// Helper function to wait for the drainer pod to start and return its logs
func waitForDrainerPod() string {
	// TODO: Implement logic to wait for drainer pod and get logs
	// Example: Use Kubernetes API to wait for drainer pod and retrieve logs
	// You may need to customize this based on your environment and deployment

	return "Drainer pod logs not implemented in the test."
}

// TODO: Add more helper functions or assertions as needed for message migration verification
// Example: assertMessageMigration(logs string)
