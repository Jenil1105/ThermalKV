from typing import List
from pydantic import BaseModel


class Analysis(BaseModel):
    update: bool
    reason: str
    sections: List[str]