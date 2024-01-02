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

	"github.com/onsi/ginkgo/v2"
	//"github.com/onsi/gomega"
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

		// Deploy the ActiveMQ Artemis broker without node affinity rules
		fileName := "ex-aao-without-affinity.yaml"
		deployBroker(fileName, dynamicClient, resourceGVR, namespace)

		// Wait for the broker pods to be created (maximum wait time: 5 minutes)
		maxWaitTime := 5 * time.Minute
		startTime := time.Now()

		for {
			if time.Since(startTime) > maxWaitTime {
				gomega.Fail("Timed out waiting for broker pods on at least two different nodes")
			}

			// Check if broker pods are on at least two different nodes
			pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: "ActiveMQArtemis=ex-aao,application=ex-aao-app",
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error getting broker pod list")

			// Check if broker pods are on at least two different nodes
			nodes := make(map[string]struct{})
			for _, pod := range pods.Items {
				nodeName := pod.Spec.NodeName
				nodes[nodeName] = struct{}{}
			}

			// Confirm that there are at least two different nodes
			if len(nodes) > 1 {
				break
			}

			// Wait for a short interval before checking again
			time.Sleep(5 * time.Second)
		}

		fmt.Println("All ActiveMQArtemis broker pods are on different nodes.")
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
