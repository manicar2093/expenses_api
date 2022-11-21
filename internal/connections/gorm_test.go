package connections_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/manicar2093/expenses_api/internal/connections"
)

var _ = Describe("Gorm", func() {

	It("creates an instance of GORM", func() {
		got := connections.GetGormConnection()

		Expect(got).ToNot(BeNil())
	})

})
