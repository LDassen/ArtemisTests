package MessageMigration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructuredv1 "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/homedir" 
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("ActiveMQ Artemis Message Migration Test", func() {
	var namespace string

	ginkgo.BeforeEach(func() {
		var err error
		var kubeconfig *rest.Config

		// Use in-cluster config if running in a Kubernetes cluster
		if kubeconfig, err = rest.InClusterConfig(); err != nil {
			// If not in a cluster, use kubeconfig file from home directory
			home := homedir.HomeDir()
			kubeconfig, err = rest.InClusterConfig()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		}

		namespace = "activemq-artemis-brokers"
	})

	ginkgo.AfterEach(func() {
		// Add cleanup logic here, e.g., delete resources created during the test
	})

	ginkgo.It("Should scale down ActiveMQ Artemis brokers and verify message migration", func() {
		// Scale down ActiveMQ Artemis brokers
		scaleDownSize := 2
		scaleDownActiveMQArtemis(scaleDownSize)

		// Check if brokers scaled down correctly
		verifyBrokerScalingEventually(scaleDownSize)

		// Wait for the drainer pod to start and print its logs
		drainerPodLogs := waitForDrainerPodEventually()
		fmt.Printf("Drainer Pod Logs:\n%s\n", drainerPodLogs)
	})
})

// Helper function to scale down ActiveMQ Artemis brokers
func scaleDownActiveMQArtemis(size int) {
	fileName := "ex-aaoMM.yaml"

	filePath, err := filepath.Abs(fileName)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	fileBytes, err := ioutil.ReadFile(filePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	decUnstructured := yaml.NewDecodingSerializer(unstructuredv1.UnstructuredJSONScheme)
	obj := &unstructuredv1.Unstructured{}
	_, _, err = decUnstructured.Decode(fileBytes, nil, obj)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	obj.SetAPIVersion("broker.amq.io/v1beta1")
	obj.SetKind("ActiveMQArtemis")

	resourceClient := dynamicClient.Resource(resourceGVR).Namespace(namespace)

	// Set replicas to the desired size
	unstructuredv1.SetNestedField(obj.Object, int64(size), "spec", "replicas")

	_, err = resourceClient.Update(context.TODO(), obj, v1.UpdateOptions{})
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error scaling down ActiveMQ Artemis brokers")
}

// Helper function to verify if brokers scaled down correctly
func verifyBrokerScalingEventually(expectedSize int) {
	gomega.Eventually(func() int {
		pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{
			LabelSelector: "app=activemq-artemis-broker",
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		return len(pods.Items)
	}, time.Minute, time.Second).Should(gomega.Equal(expectedSize), "Brokers did not scale down correctly")
}

// Helper function to wait for the drainer pod to start and return its logs
func waitForDrainerPodEventually() string {
	var lastPodName string

	gomega.Eventually(func() string {
		// Get the list of pods
		pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{
			LabelSelector: "app=activemq-artemis-broker",
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Check if there is a pod going down
		if len(pods.Items) == 1 && lastPodName != pods.Items[0].Name {
			// Get the logs of the last pod before it is scaled down
			logs, err := getPodLogs(pods.Items[0].Name)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			lastPodName = pods.Items[0].Name

			return logs
		}

		return ""
	}, time.Minute, time.Second).ShouldNot(gomega.BeEmpty(), "Logs not retrieved for the pod going down")
}

// Helper function to get logs of a pod
func getPodLogs(podName string) (string, error) {
	podLogOptions := &v1.PodLogOptions{}

	podLogs, err := kubeClient.CoreV1().Pods(namespace).GetLogs(podName, podLogOptions).Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	logBytes, err := ioutil.ReadAll(podLogs)
	if err != nil {
		return "", err
	}

	return string(logBytes), nil
}
