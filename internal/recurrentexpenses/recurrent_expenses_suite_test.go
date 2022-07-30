package recurrentexpenses_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRecurrentExpenses(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RecurrentExpenses Suite")
}
