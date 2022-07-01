
from src.expenses import models, schemas
from src.expenses.repositories.interfaces import IExpensesRepository
from src.expenses.services.interfaces import IExpensesService


class ExpensesServiceImpl(IExpensesService):

    def __init__(self, repo: IExpensesRepository) -> 'ExpensesServiceImpl':
        self.expenses_repo: IExpensesRepository = repo

    def create_expense(self, expense: schemas.CreateExpense) -> models.Expense:
        new_expense = models.Expense(
            amount=expense.amount,
            description=expense.description,
        )
        self.expenses_repo.save(new_expense)
        return new_expense
