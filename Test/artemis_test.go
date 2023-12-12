package test_test

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Artemis Broker Pods", func() {
	It("should have the correct number of 'broker' pods running", func() {
		namespace := "activemq-artemis-operator"
		expectedPodCount := 3 // Set your expected number of 'broker' pods

		// Run kubectl command to get the number of 'broker' pods
		cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "--selector=application=ex-aao-app", "--no-headers", "-o", "custom-columns=NAME:.metadata.name")
		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		Expect(err).To(BeNil(), "Error running kubectl command: %v", err)

		// Count the number of lines in the output
		actualPodCount := len(strings.Split(out.String(), "\n")) - 1

		// Assert that the actual count matches the expected count
		Expect(actualPodCount).To(Equal(expectedPodCount), "Expected %d 'broker' pods, but found %d", expectedPodCount, actualPodCount)
	})
})

func TestArtemis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Artemis Suite")
}
