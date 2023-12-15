package MultiBrokerSetup_test

import (
	"context"
	"fmt"
	"io/ioutil"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Apply Kubernetes Configuration File and Get Error Logs", func() {
	It("should apply a configuration file and retrieve error logs", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "unexisting-namespace"
		configFilePath := "ex-aao.yaml" // Change this to the path of your single configuration file

		// Check if the namespace exists
		namespaceExists, err := checkNamespaceExists(clientset, namespace)
		Expect(err).To(BeNil(), "Error checking if namespace exists: %v", err)

		// If the namespace doesn't exist, print Kubernetes logs
		if !namespaceExists {
			namespaceLogs, err := getNamespaceLogs(clientset, namespace)
			Expect(err).To(BeNil(), "Error getting namespace logs: %v", err)

			fmt.Printf("Namespace %s does not exist. Kubernetes logs:\n%s\n", namespace, namespaceLogs)
		}

		// Apply the single Kubernetes configuration file (if the namespace exists)
		if !namespaceExists {
			err = applyConfigFile(clientset, namespace, configFilePath)
			Expect(err).ToNot(BeNil(), "Expected an error applying the configuration file")
			fmt.Printf("Error applying the configuration file: %v\n", err)
		}
	})
})

// checkNamespaceExists checks if a namespace exists in the Kubernetes cluster
func checkNamespaceExists(clientset *kubernetes.Clientset, namespace string) (bool, error) {
	_, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err == nil {
		// Namespace exists
		return true, nil
	}
	if metav1.IsNotFound(err) {
		// Namespace not found
		return false, nil
	}
	// Other error
	return false, err
}

// getNamespaceLogs retrieves logs for a namespace
func getNamespaceLogs(clientset *kubernetes.Clientset, namespace string) (string, error) {
	logOptions := &corev1.PodLogOptions{}
	logs, err := clientset.CoreV1().Namespaces().GetLogs(namespace, logOptions).Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer logs.Close()

	logData, err := ioutil.ReadAll(logs)
	if err != nil {
		return "", err
	}

	return string(logData), nil
}

// applyConfigFile applies a single Kubernetes configuration file
func applyConfigFile(clientset *kubernetes.Clientset, namespace, configFilePath string) error {
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	_, err = clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(namespace).
		Body(content).
		DoRaw(context.TODO())
	if err != nil {
		return err
	}

	return nil
}