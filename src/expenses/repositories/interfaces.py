from abc import ABC, abstractmethod

from src.entities.expense import Expense


class IExpensesRepository(ABC):
    @abstractmethod
    def save(self, expense: Expense):
        raise NotImplementedError
