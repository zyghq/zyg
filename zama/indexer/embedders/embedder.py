from typing import Sequence

import ollama
from pydantic import BaseModel

from indexer.base import ChunkNode


class Embedding(BaseModel):
    embedding: Sequence[float]
    chunk_id: str | None = None
    chunk_counter: int | None = None
    content: str
    metadata: dict | None = None
    parent: str | None = None
    child: str | None = None


class LocalOllamaEmbedding:
    def __init__(self, embed_model):
        self.embed_model = embed_model

    def generate(self, chunk: ChunkNode):
        result: dict = ollama.embeddings(model=self.embed_model, prompt=chunk.content)  # type: ignore (issue with ollama's return type)
        return Embedding(
            embedding=result["embedding"],
            chunk_id=chunk.chunk_id,
            chunk_counter=chunk.chunk_counter,
            content=chunk.content,
            metadata=chunk.metadata,
            parent=chunk.parent.chunk_id if chunk.parent else None,
            child=chunk.child.chunk_id if chunk.child else None,
        )
