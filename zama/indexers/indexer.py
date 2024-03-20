from typing import List
from uuid import uuid4

from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_core.documents import Document

from core.base import WebPageContent


class ChunkNode:
    def __init__(
        self,
        *,
        chunk_id: str,
        chunk_counter: int,
        content: str,
        metadata: dict,
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


class WebDocumentSplitter:
    content_type = "text/html"

    def __init__(
        self, content: WebPageContent, *, chunk_size: int = 1024, chunk_overlap=0
    ):
        self.content = content
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap

        self.chunks = []

    def _split(self) -> List[Document]:
        splitter = RecursiveCharacterTextSplitter(
            chunk_size=self.chunk_size, chunk_overlap=self.chunk_overlap
        )
        docs = [
            Document(page_content=self.content.content, metadata=self.content.metadata)
        ]
        splits = splitter.split_documents(documents=docs)
        return splits

    def split(self) -> List[ChunkNode]:
        splits = self._split()
        prev = None
        for i, item in enumerate(splits):
            chunk = ChunkNode(
                chunk_id=str(uuid4()),
                chunk_counter=i,
                content=item.page_content,
                metadata=item.metadata,
            )
            if prev:
                prev.child = chunk
                chunk.parent = prev
            prev = chunk
            self.chunks.append(chunk)
        return self.chunks

    def to_dict(self):
        return [chunk.to_dict() for chunk in self.chunks]
