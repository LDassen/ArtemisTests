package MessageMigration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
)

var _ = ginkgo.Describe("ActiveMQ Artemis Message Migration Test", func() {
	var dynamicClient dynamic.Interface
	var kubeClient *kubernetes.Clientset
	var namespace string
	var resourceGVR schema.GroupVersionResource

	ginkgo.BeforeEach(func() {
		var err error
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		dynamicClient, err = dynamic.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		namespace = "activemq-artemis-brokers"
		resourceGVR = schema.GroupVersionResource{
			Group:    "broker.amq.io",
			Version:  "v1beta1",
			Resource: "activemqartemises",
		}
	})

	ginkgo.It("Should perform ActiveMQ Artemis Message Migration", func() {
		// Apply the YAML file to create ActiveMQArtemis resource
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

		createdObj, err := resourceClient.Create(context.TODO(), obj, metav1.CreateOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error creating ActiveMQArtemis resource")

		// Confirm that the resource has been created
		fmt.Printf("Created ActiveMQArtemis resource: %s\n", createdObj.GetName())

		// Wait for some time for the pods to stabilize
		time.Sleep(2 * time.Minute)

		// Check if the starting situation of 3 brokers becomes 2 brokers
		labelSelector := "ActiveMQArtemis=ex-aao,application=ex-aao-app"
		expectedReplicaCount := 2
		podsReady, err := arePodsReady(namespace, labelSelector, expectedReplicaCount)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error checking pod readiness")

		gomega.Expect(podsReady).To(gomega.BeTrue(), "Pods are not ready")

		// TODO: Implement logic to check the logs from the broker that goes down
		// For example, you can use kubeClient.CoreV1().Pods(namespace).GetLogs() to fetch logs.

		// TODO: Implement logic to check the logs for the drainer pod that gets started and closed
		// For example, you can use kubeClient.CoreV1().Pods(namespace).GetLogs() to fetch logs.
	})
})

func arePodsReady(namespace, labelSelector string, expectedReplicaCount int) (bool, error) {
	pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return false, err
	}

	readyCount := 0
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodRunning {
			readyCount++
		}
	}

	return read
}
