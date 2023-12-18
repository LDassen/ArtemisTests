package AutoCreationQueue_test

import (
    "bytes"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/util/remotecommand"
    "k8s.io/api/core/v1"

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
        // Assuming the configuration is set correctly in your environment
        var err error
        config, err = rest.InClusterConfig()
        Expect(err).NotTo(HaveOccurred())
        clientset, err = kubernetes.NewForConfig(config)
        Expect(err).NotTo(HaveOccurred())
    })

    It("should send messages and check queue", func() {
        // Find the Artemis broker pod
        pods, err := clientset.CoreV1().Pods("activemq-artemis-brokers").List(context.TODO(), metav1.ListOptions{
            LabelSelector: "application=ex-aao-app",
        })
        Expect(err).NotTo(HaveOccurred())
        Expect(pods.Items).NotTo(BeEmpty())

        podName := pods.Items[0].Name

        // Execute the Artemis producer command
        execProducerCmd := []string{"./amq-broker/bin/artemis", "producer", "--user", "cgi", "--password", "cgi", "--url", "tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616", "--message-count", "100"}
        _, err = execCommandInPod(clientset, config, podName, "activemq-artemis-brokers", execProducerCmd)
        Expect(err).NotTo(HaveOccurred())

        // Execute the Artemis queue stat command
        execQueueStatCmd := []string{"./amq-broker/bin/artemis", "queue", "stat", "--url", "tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616", "--user", "cgi", "--password", "cgi", "--maxRows", "200", "--clustered"}
        _, err = execCommandInPod(clientset, config, podName, "activemq-artemis-brokers", execQueueStatCmd)
        Expect(err).NotTo(HaveOccurred())
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


# some more ls aliases
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'
alias kall='kubectl get all -n kafka-brokers'
alias kalld='kubectl get all -n default'
# alias kedit='kubectl edit kafkaclusters kafka-brokers -n kafka-brokers-l'
alias kedit='kubectl edit kafkacluster kafka-brokers -n kafka-brokers'
alias kafka='cd /mnt/c/Users/laurent.dassen/Documents/Github/CMS-MCI/kafka'
alias mci='cd /mnt/c/Users/laurent.dassen/Documents/Github/mci-general/container-manifests'
alias artemis='cd /mnt/c/Users/laurent.dassen/Documents/Github/CMS-MCI/artemis'
alias tkube='cd /mnt/c/Users/laurent.dassen/Documents/Github/CMS-MCI/testkube'

# Current-context
alias contextmci='kubectl config use-context cms-mci-cluster'
alias contextrancher='kubectl config use-context rancher-desktop'
alias contextdev='kubectl config use-context cms-k8s-dev-aks001'
