from decimal import Decimal
from typing import Optional
from uuid import UUID

from pydantic import BaseModel


class CreateExpense(BaseModel):
    id: Optional[int]
    amount: Decimal
    description: Optional[str]

    class Config:
        orm_mode = True
