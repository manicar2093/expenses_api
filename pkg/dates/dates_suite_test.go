package dates_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDates(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dates Suite")
}
