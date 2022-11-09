package periodizer_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPeriodizer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Periodizer Suite")
}
