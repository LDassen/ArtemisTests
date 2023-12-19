package AutoCreationQueue_test

import (
    "context"
    "github.com/streadway/amqp"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "bytes"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    v1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/remotecommand"
)

var _ = Describe("Artemis Test", func() {
    var (
        clientset *kubernetes.Clientset
        namespace string = "activemq-artemis-brokers"
    )

    BeforeEach(func() {
        // Set up the client
        config, err := rest.InClusterConfig()
        Expect(err).NotTo(HaveOccurred())
        clientset, err = kubernetes.NewForConfig(config)
        Expect(err).NotTo(HaveOccurred())

        // Connect to Artemis using AMQP
        conn, err := amqp.Dial("amqp://cgi:cgi@ex-aao-hdls-svc.activemq-artemis-brokers.svc.cluster.local:61619")
        Expect(err).NotTo(HaveOccurred())
        defer conn.Close()

        ch, err := conn.Channel()
        Expect(err).NotTo(HaveOccurred())
        defer ch.Close()

        body := "Hi, this is a test!"
        err = ch.Publish(
            "",         // exchange
            "TESTKUBE", // routing key (queue name)
            false,      // mandatory
            false,      // immediate
            amqp.Publishing{
                ContentType: "text/plain",
                Body:        []byte(body),
            })
        Expect(err).NotTo(HaveOccurred())
    })

    It("should execute a command in one of the pods", func() {
        pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
            LabelSelector: "application=ex-aao-app",
        })
        Expect(err).NotTo(HaveOccurred())

        // Assuming you pick the first pod for simplicity
        podName := pods.Items[0].Name
        cmd := []string{"/bin/sh", "-c", "./amq-broker/bin/artemis queue stat --url tcp://ex-aao-hdls-svc.activemq-artemis-brokers:61616 --user cgi --password cgi --maxRows 200"}

        stdout, stderr, err := ExecuteRemoteCommand(podName, namespace, cmd, clientset.RESTConfig())
        Expect(err).NotTo(HaveOccurred())

        // Add your logic here to check for the specific number in the output line with 'TEST'
        // e.g., parse stdout and look for the desired pattern
    })
})

func ExecuteRemoteCommand(podName string, namespace string, command []string, restCfg *rest.Config) (string, string, error) {
    coreClient, err := kubernetes.NewForConfig(restCfg)
    if err != nil {
        return "", "", err
    }

    buf := &bytes.Buffer{}
    errBuf := &bytes.Buffer{}
    request := coreClient.CoreV1().RESTClient().
        Post().
        Namespace(namespace).
        Resource("pods").
        Name(podName).
        SubResource("exec").
        VersionedParams(&v1.PodExecOptions{
            Command: command,
            Stdout:  true,
            Stderr:  true,
        }, scheme.ParameterCodec)

    exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
    if err != nil {
        return "", "", err
    }

    err = exec.Stream(remotecommand.StreamOptions{
        Stdout: buf,
        Stderr: errBuf,
    })
    if err != nil {
        return "", "", err
    }

    return buf.String(), errBuf.String(), nil
}
