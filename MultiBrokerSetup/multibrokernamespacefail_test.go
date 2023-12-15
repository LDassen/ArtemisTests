package MultiBrokerSetup_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Apply Kubernetes Configuration Files and Check Logs", func() {
	It("should apply configuration files and check logs for errors", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "unexisting-namespace"
		configFilesPath := "ex-aao.yaml"

		// Apply Kubernetes configuration files
		err = applyConfigFiles(clientset, namespace, configFilesPath)
		Expect(err).To(BeNil(), "Error applying configuration files: %v", err)

		// Wait for pods to be ready
		err = waitForPodsReady(clientset, namespace, 1*time.Minute)
		Expect(err).To(BeNil(), "Error waiting for pods to be ready: %v", err)

		// Check Kubernetes logs for errors
		checkLogsForErrors(clientset, namespace)
	})
})

// applyConfigFiles applies Kubernetes configuration files in a specified directory
func applyConfigFiles(clientset *kubernetes.Clientset, namespace, configFilesPath string) error {
	files, err := ioutil.ReadDir(configFilesPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(configFilesPath, file.Name())
		content, err := ioutil.ReadFile(filePath)
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
	}

	return nil
}

// waitForPodsReady waits for all pods in a namespace to be in the 'Running' phase
func waitForPodsReady(clientset *kubernetes.Clientset, namespace string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		for _, pod := range pods.Items {
			if pod.Status.Phase != corev1.PodRunning {
				return false, nil
			}
		}

		return true, nil
	})
}

// checkLogsForErrors retrieves logs from all pods in a namespace and prints any errors
func checkLogsForErrors(clientset *kubernetes.Clientset, namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	Expect(err).To(BeNil(), "Error getting pods: %v", err)

	for _, pod := range pods.Items {
		logs, err := getPodLogs(clientset, namespace, pod.Name)
		Expect(err).To(BeNil(), "Error getting logs for pod %s: %v", pod.Name, err)

		// Check logs for errors and print
		if containsError(logs) {
			fmt.Printf("Error found in logs for pod %s:\n%s\n", pod.Name, logs)
		}
	}
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

// containsError checks if logs contain any error messages
func containsError(logs string) bool {
	// Implement your own logic to check for errors in the logs
	// For example, you can search for specific error patterns
	return strings.Contains(logs, "ERROR")
}