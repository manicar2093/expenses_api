from abc import ABC, abstractmethod

from src.expenses import schemas
from src.entities.expense import Expense


class IExpensesService(ABC):

    @abstractmethod
    def create_expense(self, expense: schemas.CreateExpense) -> Expense:
        raise NotImplementedError
