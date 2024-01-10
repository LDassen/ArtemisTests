package Deployment_test

import (
	"context"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/api/core/v1"
)

var _ = Describe("Check the ActiveMQ Artemis Broker PVCs", func() {
	It("should have 3 PVCs with status 'Bound'", func() {
		config, err := rest.InClusterConfig()
		Expect(err).To(BeNil(), "Error getting in-cluster config: %v", err)

		clientset, err := kubernetes.NewForConfig(config)
		Expect(err).To(BeNil(), "Error creating Kubernetes client: %v", err)

		namespace := "activemq-artemis-brokers"
		expectedPVCCount := 3
		expectedPVCStatus := v1.ClaimBound

		pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
		Expect(err).To(BeNil(), "Error getting PVCs: %v", err)

		var actualPVCCount int
		for _, pvc := range pvcs.Items {
			if pvc.Status.Phase == expectedPVCStatus {
				fmt.Printf("PVC Name: %s, Status: %s\n", pvc.Name, pvc.Status.Phase)
				actualPVCCount++
			}
		}

		Expect(actualPVCCount).To(Equal(expectedPVCCount), "Expected %d PVCs with status 'Bound', but found %d", expectedPVCCount, actualPVCCount)
	})
})
