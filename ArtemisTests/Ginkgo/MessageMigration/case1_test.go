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
	"k8s.io/client-go/util/homedir"  // Import homeDir from client-go util package

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var kubeClient *kubernetes.Clientset

func init() {
	gomega.RegisterFailHandler(ginkgo.Fail)
	format.UseStringRepresentation = false
}

func NewGomegaWithT(t types.GomegaTestingT) gomega.Gomega {
	return gomega.NewWithT(t)
}

// Helper function to wait for pods to be ready
func waitForPodsToBeReady(namespace, labelSelector string, expectedReplicaCount int) {
	var podsReady bool
	var err error

	timeout := time.Now().Add(time.Minute * 5)
	pollInterval := time.Second * 5

	for time.Now().Before(timeout) {
		podsReady, err = arePodsReady(namespace, labelSelector, expectedReplicaCount)
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error checking if pods are ready")

		if podsReady {
			break
		}

		time.Sleep(pollInterval)
	}

	gomega.Expect(podsReady).To(gomega.BeTrue(), "Pods did not become ready within the timeout period")
}

// Helper function to check if pods are ready
func arePodsReady(namespace, labelSelector string, expectedReplicaCount int) (bool, error) {
	pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return false, err
	}

	if len(pods.Items) != expectedReplicaCount {
		return false, nil
	}

	for _, pod := range pods.Items {
		if pod.Status.Phase != corev1.PodRunning {
			return false, nil
		}

		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status != corev1.ConditionTrue {
				return false, nil
			}
		}
	}

	return true, nil
}

// Helper function to check drainer image logs
func checkDrainerImageLogs(namespace, labelSelector string) {
	pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error listing pods")

	for _, pod := range pods.Items {
		// Fetch and check logs
		logs, err := kubeClient.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{}).DoRaw(context.TODO())
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error fetching pod logs")

		// TODO: Implement logic to check for the drainer image in logs
		// For example, you can use regular expressions or any other logic based on your application logs.
		fmt.Printf("Logs for pod %s:\n%s\n", pod.Name, logs)
	}
}

var _ = ginkgo.Describe("ActiveMQ Artemis Deployment Test", func() {
	var dynamicClient dynamic.Interface
	var namespace string
	var resourceGVR schema.GroupVersionResource

	ginkgo.BeforeSuite(func() {
		kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.BeforeEach(func() {
		var err error
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		dynamicClient, err = dynamic.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		namespace = "activemq-artemis-brokers"
		resourceGVR = schema.GroupVersionResource{
			Group:    "broker.amq.io",
			Version:  "v1beta1",
			Resource: "activemqartemises",
		}
	})

	ginkgo.It("Should create ActiveMQArtemis resource", func() {
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

		// Wait for pods to be ready
		waitForPodsToBeReady(namespace, "ActiveMQArtemis=ex-aao,application=ex-aao-app", 2)

		// Check drainer image logs
		checkDrainerImageLogs(namespace, "ActiveMQArtemis=ex-aao,application=ex-aao-app")
	})
})
