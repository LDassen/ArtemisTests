package MultiBrokerSetup

import (
	"context"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var _ = ginkgo.Describe("Artemis Broker Deployment", func() {
	var kubeClient *kubernetes.Clientset

	ginkgo.BeforeEach(func() {
		// Load the in-cluster or local Kubernetes config
		config, err := rest.InClusterConfig()
		if err != nil {
			// If running outside the cluster, use kubeconfig file
			home := homedir.HomeDir()
			kubeconfig := filepath.Join(home, ".kube", "config")
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
		}

		// Create a Kubernetes client
		kubeClient, err = kubernetes.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("should apply Artemis broker deployment to a non-existing namespace", func() {
		// Set the non-existing namespace
		nonExistingNamespace := "non-existing-namespace"

		// Create a Deployment object
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "artemis-broker",
				Namespace: nonExistingNamespace,
			},
			Spec: appsv1.DeploymentSpec{
				// Add your deployment spec here
				// ...
			},
		}

		// Apply the Deployment
		_, err := kubeClient.AppsV1().Deployments(nonExistingNamespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		gomega.Expect(err).To(gomega.HaveOccurred())
	})

	ginkgo.AfterEach(func() {
		// Cleanup if necessary
	})

})

func TestArtemisBrokerDeployment(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Artemis Broker Deployment Suite")
}
