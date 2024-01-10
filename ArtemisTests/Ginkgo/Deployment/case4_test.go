package Deployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Check Securityfile Existence", func() {
	It("should ensure the Securityfile 'ex-aao-prop' exists in the specified namespace", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artemis-brokers"
		securityfileName := "ex-aao-prop"

		statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), securityfileName, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("Securityfile '%s' not found in namespace '%s'\n", securityfileName, namespace)
		} else {
			fmt.Printf("Securityfile '%s' found in namespace '%s'\n", securityfileName, namespace)
		}
	})
})