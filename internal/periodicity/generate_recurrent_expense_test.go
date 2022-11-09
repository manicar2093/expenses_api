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
		ctx                                     context.Context
		service                                 *periodicity.ExpensePeriodicityServiceImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		recurrentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		recurrentExpensesMonthlyCreatedRepoMock = &mocks.RecurrentExpensesMonthlyCreatedRepo{}
		timeGetterMock = &mocks.TimeGetable{}
		periodizerExpensesGenMock = &mocks.ExpensesCountByRecurrentExpensePeriodicityGenerable{}
		ctx = context.Background()
		service = periodicity.NewExpensePeriodicityServiceImpl(
			expensesRepoMock,
			recurrentExpensesRepoMock,
			recurrentExpensesMonthlyCreatedRepoMock,
			timeGetterMock,
			periodizerExpensesGenMock,
		)
	})

	AfterEach(func() {
		T := GinkgoT()
		expensesRepoMock.AssertExpectations(T)
		recurrentExpensesRepoMock.AssertExpectations(T)
		recurrentExpensesMonthlyCreatedRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
		periodizerExpensesGenMock.AssertExpectations(T)
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
								ExpensesRelated: []primitive.ObjectID{
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

		When("recurrent expense does not contain periodicity", func() {
			It("creates it by default as monthly and update it at db", func() {
				var (
					expectedName                   = faker.Name()
					expectedLastCreationDate       = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
					expectedRecurrentExpenesID     = primitive.NewObjectID()
					expectedRecurrentExpenseAmount = faker.Latitude()
					expectedServiceCallMonth       = uint(11)
					expectedServiceCallYear        = uint(2022)
					expectedRecurrentExpense       = entities.RecurrentExpense{
						ID:               expectedRecurrentExpenesID,
						Name:             expectedName,
						Amount:           expectedRecurrentExpenseAmount,
						LastCreationDate: &expectedLastCreationDate,
					}
					expectedRecurrentExpenses = []*entities.RecurrentExpense{
						&expectedRecurrentExpense,
					}
					expectedToday              = expectedLastCreationDate.AddDate(0, 1, 0)
					expectedDay                = uint(expectedToday.Day())
					expectedMonth              = uint(expectedToday.Month())
					expectedYear               = uint(expectedToday.Year())
					expectedExpensesCreated    = uint(1)
					expectedExpensesIDsCreated = testfunc.SliceGenerator(expectedExpensesCreated, primitive.NewObjectID)
					expectedTotalExpensesPaid  = uint(0)
					expectedExpenseToSaves     = []*entities.Expense{
						{
							RecurrentExpenseID: expectedRecurrentExpenesID,
							Name:               expectedName,
							Amount:             expectedRecurrentExpenseAmount,
							Day:                expectedDay,
							Month:              expectedMonth,
							Year:               expectedYear,
							IsRecurrent:        true,
						},
					}
					expectedRecurrentExpensesMonthlyCreated = &entities.RecurrentExpensesMonthlyCreated{
						Month: expectedMonth,
						Year:  expectedYear,
						ExpensesCount: []*entities.ExpensesCount{
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								RecurrentExpense:   &expectedRecurrentExpense,
								ExpensesRelated:    expectedExpensesIDsCreated,
								TotalExpenses:      expectedExpensesCreated,
								TotalExpensesPaid:  0,
							},
						},
					}
					expectedRecurrentExpenseToUpdate = &entities.RecurrentExpense{
						ID:               expectedRecurrentExpenesID,
						Name:             expectedName,
						Amount:           expectedRecurrentExpenseAmount,
						Periodicity:      periodtypes.Monthly,
						LastCreationDate: &expectedToday,
					}
				)
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().FindByCurrentMonthAndYear(
					ctx,
					expectedServiceCallMonth,
					expectedServiceCallYear,
				).Return(nil, &repos.NotFoundError{})
				recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
					expectedRecurrentExpenses, nil,
				)
				timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
				periodizerExpensesGenMock.EXPECT().GenerateExpensesCountByRecurrentExpensePeriodicity(
					ctx,
					&expectedRecurrentExpense,
				).Return(expectedExpenseToSaves, nil)
				expensesRepoMock.EXPECT().SaveMany(ctx, expectedExpenseToSaves).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, expectedRecurrentExpensesMonthlyCreated).Return(nil)
				recurrentExpensesRepoMock.EXPECT().Update(ctx, expectedRecurrentExpenseToUpdate).Return(nil)

				got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

				Expect(err).ToNot(HaveOccurred())
				Expect(got.Month).To(Equal(expectedMonth))
				Expect(got.Year).To(Equal(expectedYear))
				Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
				Expect(got.ExpensesCount[0].RecurrentExpense).To(Equal(&expectedRecurrentExpense))
				Expect(got.ExpensesCount[0].ExpensesRelated).To(Equal(expectedExpensesIDsCreated))
				Expect(got.ExpensesCount[0].TotalExpenses).To(Equal(expectedExpensesCreated))
				Expect(got.ExpensesCount[0].TotalExpensesPaid).To(Equal(expectedTotalExpensesPaid))
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
					recurrentExpenseLastCreationDate = time.Date(2022, 10, 1, 0, 0, 0, 0, time.Local)
					expectedRecurrentExpense         = &entities.RecurrentExpense{
						ID:               primitive.NewObjectID(),
						Name:             faker.Name(),
						Amount:           faker.Latitude(),
						Description:      faker.Paragraph(),
						Periodicity:      periodtypes.Daily,
						LastCreationDate: &recurrentExpenseLastCreationDate,
					}
					expectedExpensesCreated                 = uint(30)
					expectedName                            = expectedRecurrentExpense.Name
					expectedRecurrentExpenesID              = expectedRecurrentExpense.ID
					expectedRecurrentDescription            = expectedRecurrentExpense.Description
					expectedRecurrentExpenseAmount          = expectedRecurrentExpense.Amount
					expectedPeriodicity                     = expectedRecurrentExpense.Periodicity
					expectedToday                           = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
					expectedDay                             = uint(expectedToday.Day())
					expectedMonth                           = uint(expectedToday.Month())
					expectedYear                            = uint(expectedToday.Year())
					expectedExpensesIDsCreated              = testfunc.SliceGenerator(expectedExpensesCreated, primitive.NewObjectID)
					expectedTotalExpensesPaid               = uint(0)
					expectedRecurrentExpensesMonthlyCreated = &entities.RecurrentExpensesMonthlyCreated{
						Month: expectedMonth,
						Year:  expectedYear,
						ExpensesCount: []*entities.ExpensesCount{
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								RecurrentExpense:   expectedRecurrentExpense,
								ExpensesRelated:    expectedExpensesIDsCreated,
								TotalExpenses:      expectedExpensesCreated,
								TotalExpensesPaid:  0,
							},
						},
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
				recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
					[]*entities.RecurrentExpense{expectedRecurrentExpense}, nil,
				)
				timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
				periodizerExpensesGenMock.EXPECT().GenerateExpensesCountByRecurrentExpensePeriodicity(
					ctx,
					expectedRecurrentExpense,
				).Return(expectedExpensesToSave, nil)
				expensesRepoMock.EXPECT().SaveMany(ctx, expectedExpensesToSave).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, expectedRecurrentExpensesMonthlyCreated).Return(nil)
				recurrentExpensesRepoMock.EXPECT().Update(ctx, expectedRecurrentExpenseToUpdate).Return(nil)

				got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

				Expect(err).ToNot(HaveOccurred())
				Expect(got.Month).To(Equal(expectedMonth))
				Expect(got.Year).To(Equal(expectedYear))
				Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
				Expect(got.ExpensesCount[0].RecurrentExpense).To(Equal(expectedRecurrentExpense))
				Expect(got.ExpensesCount[0].ExpensesRelated).To(Equal(expectedExpensesIDsCreated))
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
							Month:         expectedMonth,
							Year:          expectedYear,
							ExpensesCount: []*entities.ExpensesCount{},
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

	Describe("AddRecurrentExpense", func() {

		When("there is data saved", func() {
			It("adds new recurrent expense data", func() {
				// TODO: first periodizer must exists
			})
		})

		When("no data saved", func() {
			It("generates all month data", func() {

			})
		})
	})

	Describe("ToggleExpenseIsPaidStatus", func() {
		When("new status is false", func() {
			It("substacts 1 from paid expenses count", func() {

			})
		})
		When("new status is false", func() {
			It("adds 1 from paid expenses count", func() {

			})
		})
	})
})
