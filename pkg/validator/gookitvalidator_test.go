package validator_test

import (
	"net/http"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/pkg/errors"
	"github.com/manicar2093/expenses_api/pkg/validator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gookitvalidator", func() {

	var (
		api *validator.GooKitValidator
	)

	BeforeEach(func() {
		api = validator.NewGooKitValidator()
	})

	Describe("StructValidator", func() {

		It("returns a list of errors if any exists", func() {
			expectedDataToValidate := struct {
				Name string `validate:"required|min_len:7" json:"name,omitempty"`
			}{}
			got := api.ValidateStruct(&expectedDataToValidate)

			Expect(got).ToNot(BeNil())
			Expect(got.(errors.HandleableError).StatusCode()).To(Equal(http.StatusBadRequest))
		})

		When("there is any error", func() {
			It("returns nil", func() {
				expectedDataToValidate := struct {
					Name string `validate:"required|min_len:7" json:"name,omitempty"`
				}{
					Name: faker.Name(),
				}

				got := api.ValidateStruct(&expectedDataToValidate)

				Expect(got).To(BeNil())
			})
		})
	})

})
