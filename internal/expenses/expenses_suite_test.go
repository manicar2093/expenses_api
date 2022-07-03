package expenses_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExpenses(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Expenses Suite")
}
