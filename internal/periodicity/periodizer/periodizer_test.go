package periodizer_test

import (
	"context"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/periodicity/periodizer"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/periodtypes"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("CheckRecurrentExpensePeriodicity", func() {
	var (
		expensesRepoMock          *mocks.ExpensesRepository
		recurrentExpensesRepoMock *mocks.RecurrentExpenseRepo
		timeGetterMock            *mocks.TimeGetable
		ctx                       context.Context
		service                   *periodizer.PeriodicityService
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		recurrentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		timeGetterMock = &mocks.TimeGetable{}
		ctx = context.Background()
		service = periodizer.NewPeriodicityService(expensesRepoMock, recurrentExpensesRepoMock, timeGetterMock)
	})

	AfterEach(func() {
		T := GinkgoT()
		expensesRepoMock.AssertExpectations(T)
		recurrentExpensesRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
	})

	Describe("GenerateRecurrentExpensesByYearAndMonth", func() {

		When("recurrent expense does not contain periodicity", func() {
			It("creates it by default as monthly and set it in instance", func() {
				var (
					expectedName                   = faker.Name()
					expectedRecurrentExpenesID     = primitive.NewObjectID()
					expectedRecurrentExpenseAmount = faker.Latitude()
					expectedRecurrentExpense       = entities.RecurrentExpense{
						ID:               expectedRecurrentExpenesID,
						Name:             expectedName,
						Amount:           expectedRecurrentExpenseAmount,
						LastCreationDate: nil,
					}
					expectedToday           = time.Date(2022, 12, 1, 0, 0, 0, 0, time.Local)
					expectedExpensesCreated = 1
				)
				timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)

				got, err := service.GenerateExpensesCountByRecurrentExpensePeriodicity(ctx, &expectedRecurrentExpense)

				Expect(err).ToNot(HaveOccurred())
				Expect(expectedRecurrentExpense.Periodicity).To(Equal(periodtypes.Monthly))
				Expect(got).To(HaveLen(expectedExpensesCreated))
			})
		})

		When("recurrent expense has periodicity but not last creation date", func() {
			It("instantiate expenses and set last creation date", func() {
				var (
					expectedName                   = faker.Name()
					expectedRecurrentExpenesID     = primitive.NewObjectID()
					expectedRecurrentExpenseAmount = faker.Latitude()
					expectedPeriodicity            = periodtypes.BiMonthly
					expectedRecurrentExpense       = entities.RecurrentExpense{
						ID:               expectedRecurrentExpenesID,
						Name:             expectedName,
						Amount:           expectedRecurrentExpenseAmount,
						Periodicity:      expectedPeriodicity,
						LastCreationDate: nil,
					}
					expectedToday           = time.Date(2022, 12, 1, 0, 0, 0, 0, time.Local)
					expectedExpensesCreated = 1
				)
				timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)

				got, err := service.GenerateExpensesCountByRecurrentExpensePeriodicity(ctx, &expectedRecurrentExpense)

				Expect(err).ToNot(HaveOccurred())
				Expect(expectedRecurrentExpense.Periodicity).To(Equal(expectedPeriodicity))
				Expect(got).To(HaveLen(expectedExpensesCreated))
			})
		})

		DescribeTable("instantiate all expenses for supported periodicity", func(
			expectedRecurrentExpense *entities.RecurrentExpense,
			expectedExpensesCreated uint,
			recurrentExpenseLastCreationDate time.Time,
		) {
			expectedRecurrentExpense.LastCreationDate = &recurrentExpenseLastCreationDate
			var (
				expectedName                   = expectedRecurrentExpense.Name
				expectedRecurrentExpenesID     = expectedRecurrentExpense.ID
				expectedRecurrentDescription   = expectedRecurrentExpense.Description
				expectedRecurrentExpenseAmount = expectedRecurrentExpense.Amount
				expectedToday                  = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
				expectedDay                    = uint(expectedToday.Day())
				expectedMonth                  = uint(expectedToday.Month())
				expectedYear                   = uint(expectedToday.Year())
				expectedExpensesGenerated      = testfunc.SliceGenerator(expectedExpensesCreated, func() *entities.Expense {
					return &entities.Expense{
						RecurrentExpenseID: expectedRecurrentExpenesID,
						Name:               expectedName,
						Amount:             expectedRecurrentExpenseAmount,
						Description:        expectedRecurrentDescription,
						Day:                expectedDay,
						Month:              expectedMonth,
						Year:               expectedYear,
						IsRecurrent:        true,
					}
				})
			)
			timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)

			got, err := service.GenerateExpensesCountByRecurrentExpensePeriodicity(ctx, expectedRecurrentExpense)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(Equal(expectedExpensesGenerated))

		},
			Entry("daily expenses creates all taking days of the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.Daily,
				},
				uint(30),
				time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local),
			),
			Entry("weekly expenses creates four expenses in the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.Weekly,
				},
				uint(4),
				time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local),
			),
			Entry("fourteendaily expenses creates two expenses in the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.FourteenDays,
				},
				uint(2),
				time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local),
			),
			Entry("paydaily expenses creates two expenses in the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.Paydaily,
				},
				uint(2),
				time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local),
			),
			Entry("monthly expenses creates one expenses in the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.Monthly,
				},
				uint(1),
				time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local),
			),
			Entry("biMonthly expenses creates one expenses in the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.BiMonthly,
				},
				uint(1),
				time.Date(2022, 9, 1, 0, 0, 0, 0, time.Local),
			),
			Entry("four monthly expenses creates one expenses in the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.FourMonthly,
				},
				uint(1),
				time.Date(2022, 7, 1, 0, 0, 0, 0, time.Local),
			),
			Entry("six monthly expenses creates one expenses in the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.SixMonthly,
				},
				uint(1),
				time.Date(2022, 5, 1, 0, 0, 0, 0, time.Local),
			),
			Entry("yearly expenses creates one expenses in the month",
				&entities.RecurrentExpense{
					ID:          primitive.NewObjectID(),
					Name:        faker.Name(),
					Amount:      faker.Latitude(),
					Description: faker.Paragraph(),
					Periodicity: periodtypes.Yearly,
				},
				uint(1),
				time.Date(2021, 11, 1, 0, 0, 0, 0, time.Local),
			),
		)

		DescribeTable("when last creation is not totaly filled does not instantiate expenses", func(
			expectedLastCreationDate time.Time,
			expectedPeriodicity periodtypes.Periodicity,
		) {
			var (
				expectedName                   = faker.Name()
				expectedRecurrentExpenesID     = primitive.NewObjectID()
				expectedRecurrentExpenseAmount = faker.Latitude()
				expectedRecurrentExpenses      = entities.RecurrentExpense{
					ID:               expectedRecurrentExpenesID,
					Name:             expectedName,
					Amount:           expectedRecurrentExpenseAmount,
					Periodicity:      expectedPeriodicity,
					LastCreationDate: &expectedLastCreationDate,
				}
				expectedToday = expectedLastCreationDate.AddDate(0, 1, 0)
			)
			timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)

			got, err := service.GenerateExpensesCountByRecurrentExpensePeriodicity(ctx, &expectedRecurrentExpenses)

			Expect(err).ToNot(HaveOccurred())
			Expect(got).To(BeNil())

		},
			Entry("bi monthly expenses",
				time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local),
				periodtypes.BiMonthly,
			),
			Entry("bi monthly expenses with hour",
				time.Date(2022, 11, 1, 1, 0, 0, 0, time.Local),
				periodtypes.BiMonthly,
			),
			Entry("bi monthly expenses with minute",
				time.Date(2022, 11, 1, 0, 1, 0, 0, time.Local),
				periodtypes.BiMonthly,
			),
			Entry("bi monthly expenses with sec",
				time.Date(2022, 11, 1, 0, 0, 1, 0, time.Local),
				periodtypes.BiMonthly,
			),
			Entry("bi monthly expenses with nsec",
				time.Date(2022, 11, 1, 0, 0, 0, 1, time.Local),
				periodtypes.BiMonthly,
			),
			Entry("four monthly expenses",
				time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local),
				periodtypes.FourMonthly,
			),
			Entry("four monthly expenses with hour",
				time.Date(2022, 11, 1, 1, 0, 0, 0, time.Local),
				periodtypes.FourMonthly,
			),
			Entry("four monthly expenses with minute",
				time.Date(2022, 11, 1, 0, 1, 0, 0, time.Local),
				periodtypes.FourMonthly,
			),
			Entry("four monthly expenses with sec",
				time.Date(2022, 11, 1, 0, 0, 1, 0, time.Local),
				periodtypes.FourMonthly,
			),
			Entry("four monthly expenses with nsec",
				time.Date(2022, 11, 1, 0, 0, 0, 1, time.Local),
				periodtypes.FourMonthly,
			),
			Entry("six monthly expenses",
				time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local),
				periodtypes.SixMonthly,
			),
			Entry("six monthly expenses with hour",
				time.Date(2022, 11, 1, 1, 0, 0, 0, time.Local),
				periodtypes.SixMonthly,
			),
			Entry("six monthly expenses with minute",
				time.Date(2022, 11, 1, 0, 1, 0, 0, time.Local),
				periodtypes.SixMonthly,
			),
			Entry("six monthly expenses with sec",
				time.Date(2022, 11, 1, 0, 0, 1, 0, time.Local),
				periodtypes.SixMonthly,
			),
			Entry("six monthly expenses with nsec",
				time.Date(2022, 11, 1, 0, 0, 0, 1, time.Local),
				periodtypes.SixMonthly,
			),
			Entry("yearly monthly expenses",
				time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local),
				periodtypes.Yearly,
			),
			Entry("yearly monthly expenses with hour",
				time.Date(2022, 11, 1, 1, 0, 0, 0, time.Local),
				periodtypes.Yearly,
			),
			Entry("yearly monthly expenses with minute",
				time.Date(2022, 11, 1, 0, 1, 0, 0, time.Local),
				periodtypes.Yearly,
			),
			Entry("yearly monthly expenses with sec",
				time.Date(2022, 11, 1, 0, 0, 1, 0, time.Local),
				periodtypes.Yearly,
			),
			Entry("yearly monthly expenses with nsec",
				time.Date(2022, 11, 1, 0, 0, 0, 1, time.Local),
				periodtypes.Yearly,
			),
		)

	})

})
