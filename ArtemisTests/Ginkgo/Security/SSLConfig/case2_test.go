package SSLConfig_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"context"
)

var _ = Describe("ConfigMap Check in All Namespaces", func() {
    var clientset *kubernetes.Clientset

    BeforeEach(func() {
        // Set up the Kubernetes client
        config, err := rest.InClusterConfig()
        Expect(err).NotTo(HaveOccurred())

        clientset, err = kubernetes.NewForConfig(config)
        Expect(err).NotTo(HaveOccurred())
    })

    It("should have 'ca-bundle' ConfigMap in all namespaces", func() {
        // Get all namespaces
        namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
        Expect(err).NotTo(HaveOccurred())

        for _, namespace := range namespaces.Items {
            // Check for the presence of 'ca-bundle' ConfigMap in each namespace
            _, err := clientset.CoreV1().ConfigMaps(namespace.Name).Get(context.TODO(), "ca-bundle", metav1.GetOptions{})
            Expect(err).NotTo(HaveOccurred(), "Namespace: "+namespace.Name)
        }
    })
})
