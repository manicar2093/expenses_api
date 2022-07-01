from sqlalchemy.orm import Session

from src.entities.expense import Expense
from src.expenses.repositories.interfaces import IExpensesRepository


class ExpensesRepositoryImpl(IExpensesRepository):

    def __init__(self, session: Session) -> 'ExpensesRepositoryImpl':
        self.session: Session = session

    def save(self, expense: Expense):
        self.session.add(expense)
        self.session.commit()
        self.session.refresh(expense)
