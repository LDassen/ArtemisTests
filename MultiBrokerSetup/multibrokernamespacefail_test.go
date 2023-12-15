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

		// Apply the single Kubernetes configuration file
		err = applyConfigFile(clientset, namespace, configFilePath)
		Expect(err).ToNot(BeNil(), "Expected an error applying the configuration file")

		// Get error logs from all pods in the namespace
		errorLogs := getErrorLogs(clientset, namespace)
		fmt.Printf("Error logs from all pods:\n%s\n", errorLogs)
	})
})

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

// getErrorLogs retrieves error logs from all pods in a namespace
func getErrorLogs(clientset *kubernetes.Clientset, namespace string) string {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	Expect(err).To(BeNil(), "Error getting pods: %v", err)

	var errorLogs string

	for _, pod := range pods.Items {
		logs, err := getPodLogs(clientset, namespace, pod.Name)
		Expect(err).To(BeNil(), "Error getting logs for pod %s: %v", pod.Name, err)

		// Append logs to the overall errorLogs string
		errorLogs += fmt.Sprintf("Error logs for pod %s:\n%s\n", pod.Name, logs)
	}

	return errorLogs
}

// getPodLogs retrieves logs from a pod
func getPodLogs(clientset *kubernetes.Clientset, namespace, podName string) (string, error) {
	podLogOptions := &corev1.PodLogOptions{}
	podLogs, err := clientset.CoreV1().Pods(namespace).GetLogs(podName, podLogOptions).Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	logData, err := ioutil.ReadAll(podLogs)
	if err != nil {
		return "", err
	}

	return string(logData), nil
}
