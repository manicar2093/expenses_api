from decimal import Decimal
from typing import Optional
from uuid import UUID

from pydantic import BaseModel


class CreateExpense(BaseModel):
    amount: Decimal
    name: str
    description: Optional[str]
