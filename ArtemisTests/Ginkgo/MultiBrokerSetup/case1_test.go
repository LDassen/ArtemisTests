package MultiBrokerSetup_test

import (
    // ... [other imports] ...

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "path/filepath"
    "io/ioutil"

    "github.com/onsi/ginkgo/v2"
    "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("ActiveMQ Artemis Deployment Test", func() {
    var clientset *kubernetes.Clientset
    var dynamicClient dynamic.Interface
    var resourceName string
    var namespace string
    var resourceGVR schema.GroupVersionResource

    ginkgo.BeforeEach(func() {
        var err error
        config, err := rest.InClusterConfig()
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        clientset, err = kubernetes.NewForConfig(config)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        dynamicClient, err = dynamic.NewForConfig(config)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        namespace = "activemq-artemis-brokers"
        resourceName = "ex-aao"
        resourceGVR = schema.GroupVersionResource{
            Group:    "broker.amq.io",
            Version:  "v1beta1",
            Resource: "activemqartemises",
        }
    })

    ginkgo.It("Should create ActiveMQArtemis resource", func() {
        fileName := "ex-aao.yaml"

        filePath, err := filepath.Abs(fileName)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        fileBytes, err := ioutil.ReadFile(filePath)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        decUnstructured := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
        obj := &unstructured.Unstructured{}
        _, gvk, err := decUnstructured.Decode(fileBytes, nil, obj)
        gomega.Expect(err).NotTo(gomega.HaveOccurred())

        obj.SetAPIVersion("broker.amq.io/v1beta1")
        obj.SetKind("ActiveMQArtemis")

        resourceClient := dynamicClient.Resource(resourceGVR).Namespace(namespace)

        createdObj, err := resourceClient.Create(context.TODO(), obj, metav1.CreateOptions{})
        gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error creating ActiveMQArtemis resource")

        // Confirm that the resource has been created
        fmt.Printf("Created ActiveMQArtemis resource: %s\n", createdObj.GetName())
    })

    // No AfterEach block needed since we are not deleting the resource
})
