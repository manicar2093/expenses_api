from unittest.mock import Mock

from expects import expect, have_properties
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
            expected_name = self.fake.pystr()
            expected_schema_expense = schemas.CreateExpense(
                name=expected_name,
                amount=expected_amount,
                description=expected_description,
            )
            expected_model_expense = Expense(
                amount=expected_amount,
                description=expected_description,
                name=expected_name,
            )

            got = self.service.create_expense(expense=expected_schema_expense)

            expense_saved: Expense = self.expenses_repo_mock.save.call_args[0][0]
            expect(expense_saved).to(
                have_properties(**expected_schema_expense.dict()),
            )
