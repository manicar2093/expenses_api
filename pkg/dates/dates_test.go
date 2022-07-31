package dates_test

import (
	"time"

	"github.com/manicar2093/expenses_api/pkg/dates"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dates", func() {
	var (
		api *dates.TimeGetter
	)

	BeforeEach(func() {
		api = &dates.TimeGetter{}
	})

	Describe("GetNextMonthAtFirtsDay", func() {
		It("returns a time plus 1 month", func() {
			got := api.GetNextMonthAtFirtsDay()

			Expect(got.Month()).To(Equal(time.Now().Month() + 1))
			Expect(got.Day()).To(Equal(1))
		})
	})
})
