from unittest import expectedFailure
from unittest.mock import Mock, call

from expects import be_true, equal, expect
from faker import Faker
from mamba import before, describe, it

from src.entities import Expense
from src.expenses.repositories.repo import ExpensesRepositoryImpl

with describe(ExpensesRepositoryImpl) as self:
    with before.all:
        self.fake = Faker()
    with before.each:
        self.session_mock = Mock()
        self.repo = ExpensesRepositoryImpl(session=self.session_mock)

    with describe(ExpensesRepositoryImpl.save):
        with it('should call session to save an expense'):
            expected_expense = Expense(
                amount=self.fake.pydecimal(), description=self.fake.pystr(),
            )

            self.repo.save(expected_expense)
            expect(self.session_mock.add.call_args).to(
                equal(call(expected_expense)),
            )
            expect(self.session_mock.commit.called).to(be_true)
            expect(self.session_mock.refresh.call_args).to(
                equal(
                    call(expected_expense),
                ),
            )
