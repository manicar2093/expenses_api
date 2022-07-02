from unittest.mock import Mock

from expects import be_an, expect
from faker import Faker
from mamba import before, describe, it

from src.entities import Expense
from src.expenses import schemas
from src.expenses.repositories.interfaces import IExpensesRepository
from src.expenses.services.create import ExpensesServiceImpl

with describe(ExpensesServiceImpl) as self:
    with before.all:
        self.fake = Faker()
    with before.each:
        self.expenses_repo_mock = Mock(spec=IExpensesRepository)
        self.service = ExpensesServiceImpl(repo=self.expenses_repo_mock)
    with describe(ExpensesServiceImpl.create_expense):
        with it('should cast schema and call repository'):
            expected_amount = self.fake.pydecimal()
            expected_description = self.fake.pystr()
            expected_schema_expense = schemas.CreateExpense(
                amount=expected_amount,
                description=expected_description,
            )
            expected_model_expense = Expense(
                amount=expected_amount, description=expected_description,
            )

            got = self.service.create_expense(expense=expected_schema_expense)

            expect(self.expenses_repo_mock.save.call_args[0][0]).to(
                be_an(Expense),
            )
