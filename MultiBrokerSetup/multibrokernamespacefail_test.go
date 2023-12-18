package MultiBrokerSetup_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Apply Kubernetes Configuration File and Get Error Logs", func() {
	It("should apply a configuration file and retrieve error logs", func() {
		namespace := "unexisting-namespace"
		configFilePath := "./ex-aao.yaml" // Change this to the path of your single configuration file

		// Load configuration from default location or provide your own kubeconfig
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		// Create Kubernetes clientset
		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		// Read the content of the configuration file
		content, err := ioutil.ReadFile(configFilePath)
		Expect(err).To(BeNil(), "Error reading configuration file: %v", err)

		// Apply the configuration file to the namespace
		err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
			_, err = clientset.CoreV1().RESTClient().
				Post().
				Resource("pods").
				Namespace(namespace).
				Body(content).
				Do(context.TODO())
			return err
		})
		Expect(err).ToNot(BeNil(), "Expected an error applying the configuration file")

		// Retrieve and print error logs (if any)
		errorLogs := getNamespaceEvents(clientset, namespace)
		fmt.Printf("Error logs from namespace %s:\n%s\n", namespace, errorLogs)
	})
})

// getNamespaceEvents retrieves Kubernetes events for a namespace
func getNamespaceEvents(clientset *kubernetes.Clientset, namespace string) string {
	events, err := clientset.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Sprintf("Error retrieving events for namespace %s: %v", namespace, err)
	}

	var eventsBuilder strings.Builder

	for _, event := range events.Items {
		eventsBuilder.WriteString(fmt.Sprintf("Event: %s, Reason: %s, Message: %s\n", event.ObjectMeta.Name, event.Reason, event.Message))
	}

	return eventsBuilder.String()
}