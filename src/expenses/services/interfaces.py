from abc import ABC, abstractmethod

from src.expenses import models, schemas


class IExpensesService(ABC):

    @abstractmethod
    def create_expense(self, expense: schemas.CreateExpense) -> models.Expense:
        raise NotImplementedError
