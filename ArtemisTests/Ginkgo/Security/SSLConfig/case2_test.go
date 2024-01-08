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

    It("should have 'ca-bundle' ConfigMap with 'SYNCED=true' in all namespaces", func() {
        // Get all namespaces
        namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
        Expect(err).NotTo(HaveOccurred())

        for _, namespace := range namespaces.Items {
            // Retrieve the 'ca-bundle' ConfigMap in each namespace
            cm, err := clientset.CoreV1().ConfigMaps(namespace.Name).Get(context.TODO(), "ca-bundle", metav1.GetOptions{})
            Expect(err).NotTo(HaveOccurred(), "Namespace: "+namespace.Name)

            // Check if the 'SYNCED' status is 'true' in the ConfigMap annotations or data
            synced, exists := cm.Annotations["SYNCED"] // or use cm.Data based on where SYNCED is stored
            Expect(exists).To(BeTrue(), "SYNCED annotation missing in ca-bundle ConfigMap in Namespace: "+namespace.Name)
            Expect(synced).To(Equal("true"), "SYNCED is not true in ca-bundle ConfigMap in Namespace: "+namespace.Name)
        }
    })
})
