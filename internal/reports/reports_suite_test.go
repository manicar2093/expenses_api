package reports_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReports(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reports Suite")
}
