from typing import Dict

from pydantic import BaseModel


class WebPageContent(BaseModel):
    uid: str
    url: str
    metadata: Dict[str, str]
    content: str
