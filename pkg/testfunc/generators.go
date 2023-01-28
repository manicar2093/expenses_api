package testfunc

import (
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/pkg/validator"
	"gopkg.in/guregu/null.v4"
)

func GeneratePaidExpense() *entities.Expense {
	return &entities.Expense{
		ID:          uuid.New(),
		Name:        null.StringFrom(faker.Name()),
		Amount:      faker.Latitude(),
		UserID:      uuid.New(),
		Description: null.StringFrom(faker.Paragraph()),
		IsPaid:      true,
		CreatedAt:   ToPointer(time.Now()),
		UpdatedAt:   ToPointer(time.Now()),
	}
}

func CreateJsonRequestForTest(method, target string, body io.Reader) (req *http.Request) {
	req = httptest.NewRequest(method, target, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return
}

func EchoWithValidator() *echo.Echo {
	var e = echo.New()
	e.Validator = validator.NewGooKitValidator()
	return e
}

func SliceGenerator[T any](count int, generator func() T) []T {
	res := []T{}
	for i := 0; i < count; i++ {
		res = append(res, generator())
	}
	return res
}
