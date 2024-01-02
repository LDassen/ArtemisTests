package NodeAffinity_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/wait"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
)

var _ = ginkgo.Describe("ActiveMQ Artemis Node Affinity Test", func() {
	var kubeClient *kubernetes.Clientset
	var dynamicClient dynamic.Interface
	var namespace string
	var resourceGVR schema.GroupVersionResource

	ginkgo.BeforeEach(func() {
		var err error
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		kubeClient, err = kubernetes.NewForConfig(config)
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

	ginkgo.It("Should have ActiveMQArtemis pods on at least two different nodes", func() {
		// Delete the specific instance (CR) of ActiveMQArtemis
		err := dynamicClient.Resource(resourceGVR).Namespace(namespace).Delete(context.TODO(), "ex-aao", metav1.DeleteOptions{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error deleting ActiveMQArtemis resource")

		// Wait for the CR deletion to complete
		err = wait.PollImmediate(time.Second, time.Minute, func() (bool, error) {
			_, err := dynamicClient.Resource(resourceGVR).Namespace(namespace).Get(context.TODO(), "ex-aao", metav1.GetOptions{})
			return err != nil, nil
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error waiting for ActiveMQArtemis resource deletion")

		// Deploy the ActiveMQ Artemis broker without node affinity rules
		fileName := "ex-aao-without-affinity.yaml"
		deployBroker(fileName, dynamicClient, resourceGVR, namespace)

		// Wait for the broker pods to be created
		err = wait.PollImmediate(time.Second, time.Minute, func() (bool, error) {
			pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: "ActiveMQArtemis=ex-aao,application=ex-aao-app",
			})
			if err != nil {
				return false, nil
			}

			// Check if broker pods are on at least two different nodes
			nodes := make(map[string]struct{})
			for _, pod := range pods.Items {
				nodeName := pod.Spec.NodeName
				nodes[nodeName] = struct{}{}
			}

			// Confirm that there are at least two different nodes
			return len(nodes) > 1, nil
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error waiting for broker pods on at least two different nodes")
	})
})

func deployBroker(fileName string, dynamicClient dynamic.Interface, resourceGVR schema.GroupVersionResource, namespace string) {
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

	fmt.Printf("Created ActiveMQArtemis resource: %s\n", createdObj.GetName())
}
