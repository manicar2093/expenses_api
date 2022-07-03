package incomes_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIncomes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Incomes Suite")
}
