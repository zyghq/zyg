from typing import Dict

from pydantic import BaseModel


class WebPageContent(BaseModel):
    uid: str
    url: str
    metadata: Dict[str, str]
    content: str


class ChunkNode:
    def __init__(
        self,
        *,
        chunk_id: str,
        chunk_counter: int,
        content: str,
        metadata: Dict[str, str | int],
        parent: "ChunkNode | None" = None,
        child: "ChunkNode | None" = None,
    ):
        self.chunk_id = chunk_id
        self.chunk_counter = chunk_counter
        self.content = content
        self.metadata = metadata
        self.parent = parent
        self.child = child

    def to_dict(self):
        return {
            "chunk_id": self.chunk_id,
            "chunk_counter": self.chunk_counter,
            "content": self.content,
            "metadata": self.metadata,
            "parent": self.parent.chunk_id if self.parent else None,
            "child": self.child.chunk_id if self.child else None,
        }
