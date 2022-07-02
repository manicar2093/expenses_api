from src.entities import Expense
from src.expenses import schemas
from src.expenses.repositories.interfaces import IExpensesRepository
from src.expenses.services.interfaces import IExpensesService


class ExpensesServiceImpl(IExpensesService):

    def __init__(self, repo: IExpensesRepository) -> 'ExpensesServiceImpl':
        self.expenses_repo: IExpensesRepository = repo

    def create_expense(self, expense: schemas.CreateExpense) -> Expense:
        new_expense = Expense(
            amount=expense.amount,
            description=expense.description,
        )
        self.expenses_repo.save(new_expense)
        return new_expense
