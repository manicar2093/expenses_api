from abc import ABC, abstractmethod

from src.entities.expense import Expense
from src.expenses import schemas


class IExpensesService(ABC):

    @abstractmethod
    def create_expense(self, expense: schemas.CreateExpense) -> Expense:
        raise NotImplementedError
