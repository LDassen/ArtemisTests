package test_test

import (
	"os"
	"testing"
	"path/filepath"
    "net/http"
    . "github.com/onsi/ginkgo/v2"

    . "github.com/onsi/gomega"

)


var _ = Describe("URL Test", func() {

    It("should visit testkube.io and return 200", func() {

        resp, err := http.Get("https://testkube.io")

        Expect(err).To(BeNil())

        Expect(resp.StatusCode).To(Equal(200))

    })

})