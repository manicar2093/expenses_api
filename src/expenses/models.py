from sqlalchemy import DECIMAL, Column, Integer, String

from src.connections.database import Base


class Expense(Base):
    __tablename__ = 'expenses'

    id = Column(Integer, primary_key=True, index=True)
    amount = Column(DECIMAL, nullable=False)
    description = Column(String, nullable=True)
