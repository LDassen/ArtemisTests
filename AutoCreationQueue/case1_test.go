package AutoCreationQueue_test

import (
    "bytes"
    "context"
    "fmt"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "k8s.io/client-go/util/exec"
    "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/util/intstr"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("Artemis Test", func() {
    var (
        config *rest.Config
        clientset *kubernetes.Clientset
    )

    BeforeEach(func() {
        // Setup Kubernetes config and client
        var kubeconfig string
        if home := homedir.HomeDir(); home != "" {
            kubeconfig = filepath.Join(home, ".kube", "config")
        }
        config, _ = clientcmd.BuildConfigFromFlags("", kubeconfig)
        clientset, _ = kubernetes.NewForConfig(config)
    })

    It("should send messages and check queue", func() {
        // Find the Artemis broker pod
        pods, err := clientset.CoreV1().Pods("activemq-artemis-brokers").List(context.TODO(), metav1.ListOptions{
            LabelSelector: "application=ex-aao-app", // Using your provided label selector
        })
        Expect(err).NotTo(HaveOccurred())
        Expect(pods.Items).NotTo(BeEmpty())

        podName := pods.Items[0].Name // Assuming you want to use the first pod

        // Execute the Artemis producer command
        execProducerCmd := []string{"./amq-broker/bin/artemis", "producer", "--user", "cgi", "--password", "cgi", "--url", "tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616", "--message-count", "100"}
        producerOutput, err := execCommandInPod(clientset, config, podName, "activemq-artemis-brokers", execProducerCmd)
        Expect(err).NotTo(HaveOccurred())

        // Logic to verify producer command output (if needed)

        // Execute the Artemis queue stat command
        execQueueStatCmd := []string{"./amq-broker/bin/artemis", "queue", "stat", "--url", "tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616", "--user", "cgi", "--password", "cgi", "--maxRows", "200", "--clustered"}
        queueStatOutput, err := execCommandInPod(clientset, config, podName, "activemq-artemis-brokers", execQueueStatCmd)
        Expect(err).NotTo(HaveOccurred())

        // Logic to parse queueStatOutput and verify if 'TEST' queue has 300 messages
    })
})

func execCommandInPod(clientset *kubernetes.Clientset, config *rest.Config, podName, namespace string, command []string) (string, error) {
    req := clientset.CoreV1().RESTClient().
        Post().
        Resource("pods").
        Name(podName).
        Namespace(namespace).
        SubResource("exec").
        VersionedParams(&v1.PodExecOptions{
            Command:   command,
            Container: "ex-aao-ss-0", // Replace with your container name
            Stdout:    true,
            Stderr:    true,
        }, scheme.ParameterCodec)

    exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
    if err != nil {
        return "", err
    }

    var stdout, stderr bytes.Buffer
    err = exec.Stream(remotecommand.StreamOptions{
        Stdout: &stdout,
        Stderr: &stderr,
        Tty:    false,
    })

    if err != nil {
        return "", err
    }

    return stdout.String(), nil
}

func main() {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Artemis Test Suite")
}
