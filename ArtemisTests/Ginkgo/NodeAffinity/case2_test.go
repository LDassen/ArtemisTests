package NodeAffinity_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/yaml"
)

var _ = ginkgo.Describe("ActiveMQ Artemis Node Affinity Test", func() {
	var kubeClient *kubernetes.Clientset
	var namespace string

	ginkgo.BeforeEach(func() {
		var err error
		config, err := rest.InClusterConfig()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		namespace = "activemq-artemis-brokers"
	})

	ginkgo.AfterEach(func() {
		// Ensure cleanup after each test
		_ = deleteActiveMQArtemisResource("ex-aao")
	})

	ginkgo.It("Should have ActiveMQArtemis pods on at least two different nodes", func() {
		// Deploy the ActiveMQArtemis resource without node affinity
		deployActiveMQArtemisResource("ex-aao.yaml")

		// Wait for the deployment to stabilize
		gomega.Eventually(func() bool {
			return areActiveMQArtemisPodsOnDifferentNodes("ex-aao")
		}, time.Minute, time.Second).Should(gomega.BeTrue())
	})
})

func retryOnNotFound(action func() error) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() (bool, error) {
		err := action()
		if errors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
}

func deleteActiveMQArtemisResource(name string) error {
	// Get the dynamic client
	resourceClient := dynamicClient.Resource(resourceGVR).Namespace(namespace)

	// Delete the ActiveMQArtemis resource with retry on NotFound error
	err := retryOnNotFound(func() error {
		return resourceClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
	})
	if err != nil {
		return err
	}

	// Ignore NotFound errors and proceed
	return nil
}

func deployActiveMQArtemisResource(fileName string) {
	// Get the absolute path of the YAML file
	filePath, err := filepath.Abs(fileName)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Read the YAML file content
	fileBytes, err := ioutil.ReadFile(filePath)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Decode the YAML content to unstructured object
	decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	obj := &unstructured.Unstructured{}
	_, _, err = decUnstructured.Decode(fileBytes, nil, obj)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Set API version and kind
	obj.SetAPIVersion("broker.amq.io/v1beta1")
	obj.SetKind("ActiveMQArtemis")

	// Get the dynamic client
	resourceClient := dynamicClient.Resource(resourceGVR).Namespace(namespace)

	// Create the ActiveMQArtemis resource
	createdObj, err := resourceClient.Create(context.TODO(), obj, metav1.CreateOptions{})
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error creating ActiveMQArtemis resource")

	// Confirm that the resource has been created
	fmt.Printf("Created ActiveMQArtemis resource: %s\n", createdObj.GetName())
}

func areActiveMQArtemisPodsOnDifferentNodes(name string) bool {
	// Periodically check if broker pods are on different nodes
	timeout := time.After(time.Minute)
	ticker := time.Tick(time.Second)

	for {
		select {
		case <-timeout:
			// Timeout reached, the condition is not satisfied
			return false
		case <-ticker:
			// Get the list of broker pods in the namespace
			pods, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
				LabelSelector: "ActiveMQArtemis=" + name + ",application=" + name + "-app",
			})
			gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error getting broker pod list")

			// Check if broker pods are on different nodes
			nodes := make(map[string]struct{})
			for _, pod := range pods.Items {
				nodeName := pod.Spec.NodeName
				_, exists := nodes[nodeName]
				if exists {
					// Broker Pod is on the same node as another pod
					continue
				}
				nodes[nodeName] = struct{}{}
			}

			// Confirm that broker pods are on different nodes
			if len(nodes) > 1 {
				fmt.Println("All ActiveMQArtemis broker pods are on different nodes.")
				return true
			}
		}
	}
}
