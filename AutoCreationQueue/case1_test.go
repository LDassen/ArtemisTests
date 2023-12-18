package AutoCreationQueue_test

import (
	"context"
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/util/wait"
)

var _ = Describe("Artemis Broker Pods", func() {
	It("should have the correct number of 'broker' pods running", func() {
		namespace := "activemq-artemis-operator"
		expectedPodCount := 1 // Set your expected number of 'broker' pods

		cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "--selector=control-plane=controller-manager", "--output=jsonpath={.items[*].metadata.name}")
		session, err := cmd.CombinedOutput()

		Expect(err).To(BeNil(), "Error executing kubectl command: %v\nOutput:\n%s", err, session)

		// Debugging statements
		fmt.Printf("Retrieved pods in namespace %s:\n%s\n", namespace, session)

		// Split the output into pod names
		podNames := wait.SplitCommaSeparatedList(string(session))
		actualPodCount := len(podNames)

		for _, podName := range podNames {
			fmt.Printf("Pod Name: %s\n", podName)
			// Add more details as needed
		}

		Expect(actualPodCount).To(Equal(expectedPodCount), "Expected %d 'broker' pods, but found %d", expectedPodCount, actualPodCount)
	})
})
