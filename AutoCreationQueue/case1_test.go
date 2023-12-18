package AutoCreationQueue_test

import (
    "bytes"
    "context"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "k8s.io/client-go/kubernetes/scheme"
    executil "k8s.io/client-go/util/exec" // Corrected import
    "k8s.io/api/core/v1"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

// ... [rest of your code] ...

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

    exec, err := executil.NewSPDYExecutor(config, "POST", req.URL()) // Corrected usage
    if err != nil {
        return "", err
    }

    var stdout, stderr bytes.Buffer
    err = exec.Stream(executil.StreamOptions{ // Corrected usage
        Stdout: &stdout,
        Stderr: &stderr,
        Tty:    false,
    })

    if err != nil {
        return "", err
    }

    return stdout.String(), nil
}
