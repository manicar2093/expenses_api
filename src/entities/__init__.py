from sqlalchemy import Column, DateTime, Integer, Numeric, String, Text, text
from sqlalchemy.dialects.postgresql import TIMESTAMP

from src.connections.database import Base


class Expense(Base):
    __tablename__ = 'Expenses'

    id = Column(
        Integer, primary_key=True, server_default=text(
            "nextval('\"Expenses_id_seq\"'::regclass)",
        ),
    )
    amount = Column(Numeric(65, 30), nullable=False)
    name = Column(Text, nullable=False)
    description = Column(Text)
    created_at = Column(
        TIMESTAMP(precision=3), nullable=False,
        server_default=text('CURRENT_TIMESTAMP'),
    )
    updated_at = Column(TIMESTAMP(precision=3))


class Income(Base):
    __tablename__ = 'Incomes'

    id = Column(
        Integer, primary_key=True, server_default=text(
            "nextval('\"Incomes_id_seq\"'::regclass)",
        ),
    )
    amount = Column(Numeric(65, 30), nullable=False)
    name = Column(Text, nullable=False)
    description = Column(Text)
    created_at = Column(
        TIMESTAMP(precision=3), nullable=False,
        server_default=text('CURRENT_TIMESTAMP'),
    )
    updated_at = Column(TIMESTAMP(precision=3))
