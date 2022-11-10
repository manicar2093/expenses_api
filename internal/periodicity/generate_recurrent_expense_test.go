package periodicity_test

import (
	"context"
	"errors"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/manicar2093/expenses_api/internal/entities"
	"github.com/manicar2093/expenses_api/internal/periodicity"
	"github.com/manicar2093/expenses_api/internal/repos"
	"github.com/manicar2093/expenses_api/mocks"
	"github.com/manicar2093/expenses_api/pkg/periodtypes"
	"github.com/manicar2093/expenses_api/pkg/testfunc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ = Describe("CheckRecurrentExpensePeriodicity", func() {
	var (
		expensesRepoMock                        *mocks.ExpensesRepository
		recurrentExpensesRepoMock               *mocks.RecurrentExpenseRepo
		recurrentExpensesMonthlyCreatedRepoMock *mocks.RecurrentExpensesMonthlyCreatedRepo
		timeGetterMock                          *mocks.TimeGetable
		periodizerExpensesGenMock               *mocks.ExpensesCountByRecurrentExpensePeriodicityGenerable
		expensesCountRepoMock                   *mocks.ExpensesCountRepo
		ctx                                     context.Context
		service                                 *periodicity.ExpensePeriodicityServiceImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		recurrentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		recurrentExpensesMonthlyCreatedRepoMock = &mocks.RecurrentExpensesMonthlyCreatedRepo{}
		timeGetterMock = &mocks.TimeGetable{}
		periodizerExpensesGenMock = &mocks.ExpensesCountByRecurrentExpensePeriodicityGenerable{}
		expensesCountRepoMock = &mocks.ExpensesCountRepo{}
		ctx = context.Background()
		service = periodicity.NewExpensePeriodicityServiceImpl(
			expensesRepoMock,
			recurrentExpensesRepoMock,
			recurrentExpensesMonthlyCreatedRepoMock,
			timeGetterMock,
			periodizerExpensesGenMock,
			expensesCountRepoMock,
		)
	})

	AfterEach(func() {
		T := GinkgoT() //nolint:varnamelen
		expensesRepoMock.AssertExpectations(T)
		recurrentExpensesRepoMock.AssertExpectations(T)
		recurrentExpensesMonthlyCreatedRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
		periodizerExpensesGenMock.AssertExpectations(T)
		expensesCountRepoMock.AssertExpectations(T)
	})

	Describe("GenerateRecurrentExpensesByYearAndMonth", func() {

		When("get processed data returns an error", func() {
			It("returns error and finish process", func() {
				var (
					expectedMonth = uint(1)
					expectedYear  = uint(9000)
				)
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().FindByCurrentMonthAndYear(
					ctx,
					expectedMonth,
					expectedYear,
				).Return(nil, errors.New("an unexpected error"))

				got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedMonth, expectedYear)

				Expect(err).To(HaveOccurred())
				Expect(got).To(BeNil())
			})
		})

		When("data has been processed for requested year and month", func() {
			It("returns data from db", func() {
				var (
					expectedMonth                          = uint(1)
					expectedYear                           = uint(2022)
					expectedRecurrentExpenseMonthlyCreated = entities.RecurrentExpensesMonthlyCreated{
						ID:    primitive.NewObjectID(),
						Month: expectedMonth,
						Year:  expectedYear,
						ExpensesCount: []*entities.ExpensesCount{
							{
								RecurrentExpenseID: primitive.NewObjectID(),
								ExpensesRelatedIDs: []primitive.ObjectID{
									primitive.NewObjectID(),
									primitive.NewObjectID(),
									primitive.NewObjectID(),
								},
								TotalExpenses:     3,
								TotalExpensesPaid: 0,
							},
						},
					}
				)
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().FindByCurrentMonthAndYear(
					ctx,
					expectedMonth,
					expectedYear,
				).Return(&expectedRecurrentExpenseMonthlyCreated, nil)

				got, err := service.GenerateRecurrentExpensesByYearAndMonth(
					ctx,
					expectedMonth,
					expectedYear,
				)

				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal(&expectedRecurrentExpenseMonthlyCreated))
			})
		})

		When("data has not been processed", func() {

			var (
				expectedServiceCallMonth, expectedServiceCallYear uint
			)

			BeforeEach(func() {
				expectedServiceCallMonth = 11
				expectedServiceCallYear = 2022
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().FindByCurrentMonthAndYear(
					ctx,
					expectedServiceCallMonth,
					expectedServiceCallYear,
				).Return(nil, &repos.NotFoundError{})
			})

			It("instantiate all expenses and store them in db", func() {
				var (
					expectedRecurrentExpenseLastCreationDate = time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local)
					expectedRecurrentExpenesID               = primitive.NewObjectID()
					expectedName                             = faker.Name()
					expectedRecurrentExpenseAmount           = faker.Latitude()
					expectedRecurrentDescription             = faker.Paragraph()
					expectedPeriodicity                      = periodtypes.Daily
					expectedRecurrentExpense                 = &entities.RecurrentExpense{
						ID:               expectedRecurrentExpenesID,
						Name:             expectedName,
						Amount:           expectedRecurrentExpenseAmount,
						Description:      expectedRecurrentDescription,
						Periodicity:      expectedPeriodicity,
						LastCreationDate: &expectedRecurrentExpenseLastCreationDate,
					}
					expectedExpensesCreated                   = uint(30)
					expectedToday                             = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
					expectedDay                               = uint(expectedToday.Day())
					expectedMonth                             = uint(expectedToday.Month())
					expectedYear                              = uint(expectedToday.Year())
					expectedExpensesIDsCreated                = testfunc.SliceGenerator(expectedExpensesCreated, primitive.NewObjectID)
					expectedTotalExpensesPaid                 = uint(0)
					expectedRecurrentExpensesMonthlyCreatedID = primitive.NewObjectID()
					expectedExpensesCount1                    = entities.ExpensesCount{
						RecurrentExpenseID:                expectedRecurrentExpenesID,
						RecurrentExpensesMonthlyCreatedID: expectedRecurrentExpensesMonthlyCreatedID,
						RecurrentExpense:                  expectedRecurrentExpense,
						ExpensesRelatedIDs:                expectedExpensesIDsCreated,
						TotalExpenses:                     expectedExpensesCreated,
						TotalExpensesPaid:                 0,
					}
					expectedRecurrentExpensesMonthlyCreated = entities.RecurrentExpensesMonthlyCreated{
						Month: expectedMonth,
						Year:  expectedYear,
					}
					expectedRecurrentExpenseToUpdate = &entities.RecurrentExpense{
						ID:               expectedRecurrentExpenesID,
						Name:             expectedName,
						Amount:           expectedRecurrentExpenseAmount,
						Description:      expectedRecurrentDescription,
						Periodicity:      expectedPeriodicity,
						LastCreationDate: &expectedToday,
					}
					expectedExpensesToSave = testfunc.SliceGenerator(expectedExpensesCreated, func() *entities.Expense {
						return &entities.Expense{
							RecurrentExpenseID: &expectedRecurrentExpenesID,
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
				recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
					[]*entities.RecurrentExpense{expectedRecurrentExpense}, nil,
				)
				timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
				periodizerExpensesGenMock.EXPECT().GenerateExpensesCountByRecurrentExpensePeriodicity(
					ctx,
					expectedRecurrentExpense,
				).Return(expectedExpensesToSave, nil)
				expensesRepoMock.EXPECT().SaveMany(
					ctx,
					expectedExpensesToSave,
				).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(
					ctx,
					&expectedRecurrentExpensesMonthlyCreated,
				).Return(nil).Run(func(ctx context.Context, recurrentExpense *entities.RecurrentExpensesMonthlyCreated) {
					recurrentExpense.ID = expectedRecurrentExpensesMonthlyCreatedID
				})
				expensesCountRepoMock.EXPECT().Save(ctx, &expectedExpensesCount1).Return(nil).Once()
				recurrentExpensesRepoMock.EXPECT().Update(ctx, expectedRecurrentExpenseToUpdate).Return(nil)

				got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

				Expect(err).ToNot(HaveOccurred())
				Expect(got.Month).To(Equal(expectedMonth))
				Expect(got.Year).To(Equal(expectedYear))
				Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
				Expect(got.ExpensesCount[0].RecurrentExpense).To(Equal(expectedRecurrentExpense))
				Expect(got.ExpensesCount[0].ExpensesRelatedIDs).To(Equal(expectedExpensesIDsCreated))
				Expect(got.ExpensesCount[0].TotalExpenses).To(Equal(expectedExpensesCreated))
				Expect(got.ExpensesCount[0].TotalExpensesPaid).To(Equal(expectedTotalExpensesPaid))
			})

			When("there are not instance of expenses to be created", func() {
				It("does not register expenses", func() {
					var (
						expectedLastCreationDate       = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
						expectedPeriodicity            = periodtypes.BiMonthly
						expectedName                   = faker.Name()
						expectedRecurrentExpenesID     = primitive.NewObjectID()
						expectedRecurrentExpenseAmount = faker.Latitude()
						expectedRecurrentExpense       = entities.RecurrentExpense{
							ID:               expectedRecurrentExpenesID,
							Name:             expectedName,
							Amount:           expectedRecurrentExpenseAmount,
							Periodicity:      expectedPeriodicity,
							LastCreationDate: &expectedLastCreationDate,
						}
						expectedRecurrentExpenses = []*entities.RecurrentExpense{
							&expectedRecurrentExpense,
						}
						expectedToday                          = expectedLastCreationDate.AddDate(0, 1, 0)
						expectedMonth                          = uint(expectedToday.Month())
						expectedYear                           = uint(expectedToday.Year())
						expectedRecurrentExpenseMonthlyCreated = &entities.RecurrentExpensesMonthlyCreated{
							Month: expectedMonth,
							Year:  expectedYear,
						}
					)
					recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
						expectedRecurrentExpenses, nil,
					)
					timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
					periodizerExpensesGenMock.EXPECT().GenerateExpensesCountByRecurrentExpensePeriodicity(
						ctx,
						&expectedRecurrentExpense,
					).Return(nil, nil)
					recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, expectedRecurrentExpenseMonthlyCreated).Return(nil)

					got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

					Expect(err).ToNot(HaveOccurred())
					Expect(got.Month).To(Equal(expectedMonth))
					Expect(got.Year).To(Equal(expectedYear))
					Expect(got.ExpensesCount).To(HaveLen(0))
				})
			})

		})

	})

})
