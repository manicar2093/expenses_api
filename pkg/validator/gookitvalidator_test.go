package validator_test

import (
	"net/http"

	"github.com/bxcodec/faker/v3"
	"github.com/gookit/validate"
	"github.com/manicar2093/expenses_api/pkg/apperrors"
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
			Expect(got.(apperrors.HandleableError).StatusCode()).To(Equal(http.StatusBadRequest))
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

		When("uuid is not valid", func() {
			It("shows message for invalid uuid", func() {
				var (
					expectedSubstring      = "is not valid for UUID type"
					expectedNotUUID        = "not a uuid"
					expectedDataToValidate = struct {
						Name     string `validate:"isUUID" json:"name,omitempty"`
						LastName string `validate:"uuid" json:"last_name,omitempty"`
					}{
						Name:     expectedNotUUID,
						LastName: expectedNotUUID,
					}
				)

				got := api.ValidateStruct(&expectedDataToValidate)

				err := got.(*validator.ValidationError)
				errMap := err.Errors.(validate.Errors)
				Expect(errMap["name"]["isUUID"]).To(ContainSubstring(expectedSubstring))
				Expect(errMap["last_name"]["uuid"]).To(ContainSubstring(expectedSubstring))
			})
		})

		When("filed is required", func() {
			It("shows message for required data", func() {
				var (
					expectedSubstring      = "needs to be on request"
					expectedDataToValidate = struct {
						Name string `validate:"required" json:"name,omitempty"`
					}{}
				)

				got := api.ValidateStruct(&expectedDataToValidate)

				err := got.(*validator.ValidationError)
				errMap := err.Errors.(validate.Errors)
				Expect(errMap["name"]["required"]).To(ContainSubstring(expectedSubstring))
			})
		})
	})

})
