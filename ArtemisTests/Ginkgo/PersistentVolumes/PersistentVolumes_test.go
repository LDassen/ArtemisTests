package PersistentVolumes_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMultiBrokerSetup(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistentVolumes Suite")
}
