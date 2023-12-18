package AutoCreationQueue_test

import (
	"fmt"
	"os/exec"
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Artemis Broker Pods", func() {
	It("should have the correct number of 'broker' pods running", func() {
		namespace := "activemq-artemis-brokers"
		expectedPodCount := 3 // Set your expected number of 'broker' pods

		cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "-l", "application=ex-aao-app", "--output=json")
		output, err := cmd.CombinedOutput()
		Expect(err).To(BeNil(), "Error running kubectl command: %v\nOutput: %s", err, output)

		// Parse JSON output to get pod count
		podCount := 0
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "\"kind\": \"Pod\"") {
				podCount++
			}
		}

		// Debugging statements
		fmt.Printf("Retrieved %d pods in namespace %s\n", podCount, namespace)
		fmt.Printf("Output:\n%s\n", output)

		Expect(podCount).To(Equal(expectedPodCount), "Expected %d 'broker' pods, but found %d", expectedPodCount, podCount)
	})
})
