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
		ctx                                     context.Context
		service                                 *periodicity.ExpensePeriodicityServiceImpl
	)

	BeforeEach(func() {
		expensesRepoMock = &mocks.ExpensesRepository{}
		recurrentExpensesRepoMock = &mocks.RecurrentExpenseRepo{}
		recurrentExpensesMonthlyCreatedRepoMock = &mocks.RecurrentExpensesMonthlyCreatedRepo{}
		timeGetterMock = &mocks.TimeGetable{}
		ctx = context.Background()
		service = periodicity.NewExpensePeriodicityServiceImpl(
			expensesRepoMock,
			recurrentExpensesRepoMock,
			recurrentExpensesMonthlyCreatedRepoMock,
			timeGetterMock,
		)
	})

	AfterEach(func() {
		T := GinkgoT()
		expensesRepoMock.AssertExpectations(T)
		recurrentExpensesRepoMock.AssertExpectations(T)
		recurrentExpensesMonthlyCreatedRepoMock.AssertExpectations(T)
		timeGetterMock.AssertExpectations(T)
	})

	Describe("GenerateRecurrentExpensesByYearAndMonth", func() {

		When("get processed data returns an error", func() {
			It("returns error and finish process", func() {
				var (
					expectedMonth = uint(1)
					expectedYear  = uint(9000)
				)
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().FindByMonthAndYear(
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
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().FindByMonthAndYear(
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
				expectedServiceCallMonth = 12
				expectedServiceCallYear = 2022
				recurrentExpensesMonthlyCreatedRepoMock.EXPECT().FindByMonthAndYear(
					ctx,
					expectedServiceCallMonth,
					expectedServiceCallYear,
				).Return(nil, &repos.NotFoundError{})
			})

			Describe("create all expenses and store data in db", func() {

				Context("generates daily expenses", func() {
					It("creates all taking days of the month", func() {
						var (
							expectedName                   = faker.Name()
							expectedLastCreationDate       = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
							expectedRecurrentExpenesID     = primitive.NewObjectID()
							expectedRecurrentExpenseAmount = faker.Latitude()
							expectedRecurrentExpenses      = []*entities.RecurrentExpense{
								{
									ID:               expectedRecurrentExpenesID,
									Name:             expectedName,
									Amount:           expectedRecurrentExpenseAmount,
									Periodicity:      periodtypes.Daily,
									LastCreationDate: &expectedLastCreationDate,
								},
							}
							expectedToday              = expectedLastCreationDate.AddDate(0, 1, 0)
							expectedDay                = uint(expectedToday.Day())
							expectedMonth              = uint(expectedToday.Month())
							expectedYear               = uint(expectedToday.Year())
							expectedExpensesCreated    = uint(31)
							expectedExpensesIDsCreated = testfunc.GeneratePrimitiveObjectIDs(expectedExpensesCreated)
							expectedTotalExpensesPaid  = uint(0)
						)
						recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
							expectedRecurrentExpenses, nil,
						)
						timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
						expensesRepoMock.EXPECT().SaveMany(ctx, []*entities.Expense{
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
						}).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
						recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, &entities.RecurrentExpensesMonthlyCreated{
							Month: expectedMonth,
							Year:  expectedYear,
							ExpensesCount: []*entities.ExpensesCount{
								{
									RecurrentExpenseID: expectedRecurrentExpenesID,
									ExpensesRelated:    expectedExpensesIDsCreated,
									TotalExpenses:      expectedExpensesCreated,
									TotalExpensesPaid:  0,
								},
							},
						}).Return(nil)
						recurrentExpensesRepoMock.EXPECT().Update(ctx, &entities.RecurrentExpense{
							ID:               expectedRecurrentExpenesID,
							Name:             expectedName,
							Amount:           expectedRecurrentExpenseAmount,
							Periodicity:      periodtypes.Daily,
							LastCreationDate: &expectedToday,
						}).Return(nil)

						got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

						Expect(err).ToNot(HaveOccurred())
						Expect(got.Month).To(Equal(expectedMonth))
						Expect(got.Year).To(Equal(expectedYear))
						Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
						Expect(got.ExpensesCount[0].ExpensesRelated).To(Equal(expectedExpensesIDsCreated))
						Expect(got.ExpensesCount[0].TotalExpenses).To(Equal(expectedExpensesCreated))
						Expect(got.ExpensesCount[0].TotalExpensesPaid).To(Equal(expectedTotalExpensesPaid))
					})
				})

				Context("generates weekly expenses", func() {
					It("generates four expenses per month", func() {
						var (
							expectedName                   = faker.Name()
							expectedRecurrentExpenesID     = primitive.NewObjectID()
							expectedRecurrentExpenseAmount = faker.Latitude()
							expectedRecurrentDescription   = faker.Paragraph()
							expectedLastCreationDate       = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
							expectedRecurrentExpenses      = []*entities.RecurrentExpense{
								{
									ID:               expectedRecurrentExpenesID,
									Name:             expectedName,
									Amount:           expectedRecurrentExpenseAmount,
									Description:      expectedRecurrentDescription,
									Periodicity:      periodtypes.Weekly,
									LastCreationDate: &expectedLastCreationDate,
								},
							}

							expectedToday              = expectedLastCreationDate.AddDate(0, 1, 0)
							expectedDay                = uint(expectedToday.Day())
							expectedMonth              = uint(expectedToday.Month())
							expectedYear               = uint(expectedToday.Year())
							expectedExpensesCreated    = uint(4)
							expectedExpensesIDsCreated = testfunc.GeneratePrimitiveObjectIDs(expectedExpensesCreated)
							expectedTotalExpensesPaid  = uint(0)
						)
						recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
							expectedRecurrentExpenses, nil,
						)
						timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
						expensesRepoMock.EXPECT().SaveMany(ctx, []*entities.Expense{
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Description:        expectedRecurrentDescription,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Description:        expectedRecurrentDescription,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Description:        expectedRecurrentDescription,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Description:        expectedRecurrentDescription,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
						}).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
						recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, &entities.RecurrentExpensesMonthlyCreated{
							Month: expectedMonth,
							Year:  expectedYear,
							ExpensesCount: []*entities.ExpensesCount{
								{
									RecurrentExpenseID: expectedRecurrentExpenesID,
									ExpensesRelated:    expectedExpensesIDsCreated,
									TotalExpenses:      expectedExpensesCreated,
									TotalExpensesPaid:  0,
								},
							},
						}).Return(nil)
						recurrentExpensesRepoMock.EXPECT().Update(ctx, &entities.RecurrentExpense{
							ID:               expectedRecurrentExpenesID,
							Name:             expectedName,
							Amount:           expectedRecurrentExpenseAmount,
							Description:      expectedRecurrentDescription,
							Periodicity:      periodtypes.Weekly,
							LastCreationDate: &expectedToday,
						}).Return(nil)

						got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

						Expect(err).ToNot(HaveOccurred())
						Expect(got.Month).To(Equal(expectedMonth))
						Expect(got.Year).To(Equal(expectedYear))
						Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
						Expect(got.ExpensesCount[0].ExpensesRelated).To(Equal(expectedExpensesIDsCreated))
						Expect(got.ExpensesCount[0].TotalExpenses).To(Equal(expectedExpensesCreated))
						Expect(got.ExpensesCount[0].TotalExpensesPaid).To(Equal(expectedTotalExpensesPaid))
					})
				})

				Context("generates fourteendaily", func() {
					It("generates two expenses per month", func() {
						var (
							expectedName                   = faker.Name()
							expectedLastCreationDate       = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
							expectedRecurrentExpenesID     = primitive.NewObjectID()
							expectedRecurrentExpenseAmount = faker.Latitude()
							expectedRecurrentExpenses      = []*entities.RecurrentExpense{
								{
									ID:               expectedRecurrentExpenesID,
									Name:             expectedName,
									Amount:           expectedRecurrentExpenseAmount,
									Periodicity:      periodtypes.FourteenDays,
									LastCreationDate: &expectedLastCreationDate,
								},
							}
							expectedToday              = expectedLastCreationDate.AddDate(0, 1, 0)
							expectedDay                = uint(expectedToday.Day())
							expectedMonth              = uint(expectedToday.Month())
							expectedYear               = uint(expectedToday.Year())
							expectedExpensesCreated    = uint(2)
							expectedExpensesIDsCreated = testfunc.GeneratePrimitiveObjectIDs(expectedExpensesCreated)
							expectedTotalExpensesPaid  = uint(0)
						)
						recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
							expectedRecurrentExpenses, nil,
						)
						timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
						expensesRepoMock.EXPECT().SaveMany(ctx, []*entities.Expense{
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
						}).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
						recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, &entities.RecurrentExpensesMonthlyCreated{
							Month: expectedMonth,
							Year:  expectedYear,
							ExpensesCount: []*entities.ExpensesCount{
								{
									RecurrentExpenseID: expectedRecurrentExpenesID,
									ExpensesRelated:    expectedExpensesIDsCreated,
									TotalExpenses:      expectedExpensesCreated,
									TotalExpensesPaid:  0,
								},
							},
						}).Return(nil)
						recurrentExpensesRepoMock.EXPECT().Update(ctx, &entities.RecurrentExpense{
							ID:               expectedRecurrentExpenesID,
							Name:             expectedName,
							Amount:           expectedRecurrentExpenseAmount,
							Periodicity:      periodtypes.FourteenDays,
							LastCreationDate: &expectedToday,
						}).Return(nil)

						got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

						Expect(err).ToNot(HaveOccurred())
						Expect(got.Month).To(Equal(expectedMonth))
						Expect(got.Year).To(Equal(expectedYear))
						Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
						Expect(got.ExpensesCount[0].ExpensesRelated).To(Equal(expectedExpensesIDsCreated))
						Expect(got.ExpensesCount[0].TotalExpenses).To(Equal(expectedExpensesCreated))
						Expect(got.ExpensesCount[0].TotalExpensesPaid).To(Equal(expectedTotalExpensesPaid))
					})
				})

				Context("generates paydaily", func() {
					It("generates two expenses per month", func() {
						var (
							expectedName                   = faker.Name()
							expectedLastCreationDate       = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
							expectedRecurrentExpenesID     = primitive.NewObjectID()
							expectedRecurrentExpenseAmount = faker.Latitude()
							expectedRecurrentExpenses      = []*entities.RecurrentExpense{
								{
									ID:               expectedRecurrentExpenesID,
									Name:             expectedName,
									Amount:           expectedRecurrentExpenseAmount,
									Periodicity:      periodtypes.Paydaily,
									LastCreationDate: &expectedLastCreationDate,
								},
							}
							expectedToday              = expectedLastCreationDate.AddDate(0, 1, 0)
							expectedDay                = uint(expectedToday.Day())
							expectedMonth              = uint(expectedToday.Month())
							expectedYear               = uint(expectedToday.Year())
							expectedExpensesCreated    = uint(2)
							expectedExpensesIDsCreated = testfunc.GeneratePrimitiveObjectIDs(expectedExpensesCreated)
							expectedTotalExpensesPaid  = uint(0)
						)
						recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
							expectedRecurrentExpenses, nil,
						)
						timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
						expensesRepoMock.EXPECT().SaveMany(ctx, []*entities.Expense{
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
						}).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
						recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, &entities.RecurrentExpensesMonthlyCreated{
							Month: expectedMonth,
							Year:  expectedYear,
							ExpensesCount: []*entities.ExpensesCount{
								{
									RecurrentExpenseID: expectedRecurrentExpenesID,
									ExpensesRelated:    expectedExpensesIDsCreated,
									TotalExpenses:      expectedExpensesCreated,
									TotalExpensesPaid:  0,
								},
							},
						}).Return(nil)
						recurrentExpensesRepoMock.EXPECT().Update(ctx, &entities.RecurrentExpense{
							ID:               expectedRecurrentExpenesID,
							Name:             expectedName,
							Amount:           expectedRecurrentExpenseAmount,
							Periodicity:      periodtypes.Paydaily,
							LastCreationDate: &expectedToday,
						}).Return(nil)

						got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

						Expect(err).ToNot(HaveOccurred())
						Expect(got.Month).To(Equal(expectedMonth))
						Expect(got.Year).To(Equal(expectedYear))
						Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
						Expect(got.ExpensesCount[0].ExpensesRelated).To(Equal(expectedExpensesIDsCreated))
						Expect(got.ExpensesCount[0].TotalExpenses).To(Equal(expectedExpensesCreated))
						Expect(got.ExpensesCount[0].TotalExpensesPaid).To(Equal(expectedTotalExpensesPaid))
					})
				})

				Context("generates monthly expenses", func() {
					It("creates one expense per month", func() {
						var (
							expectedName                   = faker.Name()
							expectedLastCreationDate       = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
							expectedRecurrentExpenesID     = primitive.NewObjectID()
							expectedRecurrentExpenseAmount = faker.Latitude()
							expectedRecurrentExpenses      = []*entities.RecurrentExpense{
								{
									ID:               expectedRecurrentExpenesID,
									Name:             expectedName,
									Amount:           expectedRecurrentExpenseAmount,
									Periodicity:      periodtypes.Monthly,
									LastCreationDate: &expectedLastCreationDate,
								},
							}
							expectedToday              = expectedLastCreationDate.AddDate(0, 1, 0)
							expectedDay                = uint(expectedToday.Day())
							expectedMonth              = uint(expectedToday.Month())
							expectedYear               = uint(expectedToday.Year())
							expectedExpensesCreated    = uint(1)
							expectedExpensesIDsCreated = testfunc.GeneratePrimitiveObjectIDs(expectedExpensesCreated)
							expectedTotalExpensesPaid  = uint(0)
						)
						recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
							expectedRecurrentExpenses, nil,
						)
						timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
						expensesRepoMock.EXPECT().SaveMany(ctx, []*entities.Expense{
							{
								RecurrentExpenseID: expectedRecurrentExpenesID,
								Name:               expectedName,
								Amount:             expectedRecurrentExpenseAmount,
								Day:                expectedDay,
								Month:              expectedMonth,
								Year:               expectedYear,
								IsRecurrent:        true,
							},
						}).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
						recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, &entities.RecurrentExpensesMonthlyCreated{
							Month: expectedMonth,
							Year:  expectedYear,
							ExpensesCount: []*entities.ExpensesCount{
								{
									RecurrentExpenseID: expectedRecurrentExpenesID,
									ExpensesRelated:    expectedExpensesIDsCreated,
									TotalExpenses:      expectedExpensesCreated,
									TotalExpensesPaid:  0,
								},
							},
						}).Return(nil)
						recurrentExpensesRepoMock.EXPECT().Update(ctx, &entities.RecurrentExpense{
							ID:               expectedRecurrentExpenesID,
							Name:             expectedName,
							Amount:           expectedRecurrentExpenseAmount,
							Periodicity:      periodtypes.Monthly,
							LastCreationDate: &expectedToday,
						}).Return(nil)

						got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

						Expect(err).ToNot(HaveOccurred())
						Expect(got.Month).To(Equal(expectedMonth))
						Expect(got.Year).To(Equal(expectedYear))
						Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
						Expect(got.ExpensesCount[0].ExpensesRelated).To(Equal(expectedExpensesIDsCreated))
						Expect(got.ExpensesCount[0].TotalExpenses).To(Equal(expectedExpensesCreated))
						Expect(got.ExpensesCount[0].TotalExpensesPaid).To(Equal(expectedTotalExpensesPaid))
					})

					When("recurrent expense does not contain periodicity", func() {
						It("creates it by default as monthly and update it at db", func() {
							var (
								expectedName                   = faker.Name()
								expectedLastCreationDate       = time.Date(2022, 11, 1, 0, 0, 0, 0, time.Local)
								expectedRecurrentExpenesID     = primitive.NewObjectID()
								expectedRecurrentExpenseAmount = faker.Latitude()
								expectedRecurrentExpenses      = []*entities.RecurrentExpense{
									{
										ID:               expectedRecurrentExpenesID,
										Name:             expectedName,
										Amount:           expectedRecurrentExpenseAmount,
										LastCreationDate: &expectedLastCreationDate,
									},
								}
								expectedToday              = expectedLastCreationDate.AddDate(0, 1, 0)
								expectedDay                = uint(expectedToday.Day())
								expectedMonth              = uint(expectedToday.Month())
								expectedYear               = uint(expectedToday.Year())
								expectedExpensesCreated    = uint(1)
								expectedExpensesIDsCreated = testfunc.GeneratePrimitiveObjectIDs(expectedExpensesCreated)
								expectedTotalExpensesPaid  = uint(0)
							)
							recurrentExpensesRepoMock.EXPECT().FindAll(ctx).Return(
								expectedRecurrentExpenses, nil,
							)
							timeGetterMock.EXPECT().GetCurrentTime().Return(expectedToday)
							expensesRepoMock.EXPECT().SaveMany(ctx, []*entities.Expense{
								{
									RecurrentExpenseID: expectedRecurrentExpenesID,
									Name:               expectedName,
									Amount:             expectedRecurrentExpenseAmount,
									Day:                expectedDay,
									Month:              expectedMonth,
									Year:               expectedYear,
									IsRecurrent:        true,
								},
							}).Return(&repos.InsertManyResult{InsertedIDs: expectedExpensesIDsCreated}, nil)
							recurrentExpensesMonthlyCreatedRepoMock.EXPECT().Save(ctx, &entities.RecurrentExpensesMonthlyCreated{
								Month: expectedMonth,
								Year:  expectedYear,
								ExpensesCount: []*entities.ExpensesCount{
									{
										RecurrentExpenseID: expectedRecurrentExpenesID,
										ExpensesRelated:    expectedExpensesIDsCreated,
										TotalExpenses:      expectedExpensesCreated,
										TotalExpensesPaid:  0,
									},
								},
							}).Return(nil)
							recurrentExpensesRepoMock.EXPECT().Update(ctx, &entities.RecurrentExpense{
								ID:               expectedRecurrentExpenesID,
								Name:             expectedName,
								Amount:           expectedRecurrentExpenseAmount,
								Periodicity:      periodtypes.Monthly,
								LastCreationDate: &expectedToday,
							}).Return(nil)

							got, err := service.GenerateRecurrentExpensesByYearAndMonth(ctx, expectedServiceCallMonth, expectedServiceCallYear)

							Expect(err).ToNot(HaveOccurred())
							Expect(got.Month).To(Equal(expectedMonth))
							Expect(got.Year).To(Equal(expectedYear))
							Expect(got.ExpensesCount[0].RecurrentExpenseID).To(Equal(expectedRecurrentExpenesID))
							Expect(got.ExpensesCount[0].ExpensesRelated).To(Equal(expectedExpensesIDsCreated))
							Expect(got.ExpensesCount[0].TotalExpenses).To(Equal(expectedExpensesCreated))
							Expect(got.ExpensesCount[0].TotalExpensesPaid).To(Equal(expectedTotalExpensesPaid))
						})
					})
				})

			})

		})
	})
})
