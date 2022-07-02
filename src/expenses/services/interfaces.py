from abc import ABC, abstractmethod

from src.entities import Expense
from src.expenses import schemas


class IExpensesService(ABC):

    @abstractmethod
    def create_expense(self, expense: schemas.CreateExpense) -> Expense:
        raise NotImplementedError
