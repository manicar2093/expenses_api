package periodicity_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPeriodicity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Periodicity Suite")
}
