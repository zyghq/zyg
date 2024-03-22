from typing import Sequence

import ollama
import tiktoken
from pydantic import BaseModel

from indexer.base import ChunkNode


# taken from
# https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
def num_tokens_from_string(string: str, encoding_name: str) -> int:
    """Returns the number of tokens in a text string."""
    encoding = tiktoken.get_encoding(encoding_name)
    num_tokens = len(encoding.encode(string))
    return num_tokens


class Embedding(BaseModel):
    embed_model: str
    embedding: Sequence[float]
    approx_tokens: int
    document_uid: str
    document_url: str
    chunk_id: str
    chunk_counter: int
    content: str
    parent: str | None = None
    child: str | None = None

    def to_dict(self):
        return self.model_dump()


class LocalOllamaEmbedding:
    def __init__(self, embed_model):
        self.embed_model = embed_model

    def __call__(self, chunk: ChunkNode):
        # TODO: modify this, for handing different embed_models.
        result: dict = ollama.embeddings(model=self.embed_model, prompt=chunk.content)  # type: ignore (issue with ollama's return type)
        return Embedding(
            embed_model=self.embed_model,
            embedding=result["embedding"],
            approx_tokens=num_tokens_from_string(chunk.content, "cl100k_base"),
            document_uid=chunk.uid,
            document_url=chunk.uri,
            chunk_id=chunk.chunk_id,
            chunk_counter=chunk.chunk_counter,
            content=chunk.content,
            parent=chunk.parent.chunk_id if chunk.parent else None,
            child=chunk.child.chunk_id if chunk.child else None,
        )
