from fastapi import Depends, FastAPI
from sqlalchemy.orm import Session

from src.connections.database import create_db, get_db
from src.expenses import schemas
from src.expenses.repositories.repo import ExpensesRepositoryImpl
from src.expenses.services.create import ExpensesServiceImpl

app = FastAPI()
create_db()


@app.get('/ping')
def root():
    return 'pong'


@app.post('/expense')
def create_expense(
    expense: schemas.CreateExpense,
    db: Session = Depends(get_db),
):
    return ExpensesServiceImpl(
        repo=ExpensesRepositoryImpl(
            session=db,
        ),
    ).create_expense(expense=expense)
