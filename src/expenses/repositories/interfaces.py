from abc import ABC, abstractmethod

from src.expenses import models


class IExpensesRepository(ABC):
    @abstractmethod
    def save(self, expense: models.Expense):
        raise NotImplementedError
