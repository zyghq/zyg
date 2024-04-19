import logging
from typing import List
from uuid import uuid4

from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain_core.documents import Document

from indexer.base import ChunkNode, WebPageContent

logger = logging.getLogger(__name__)


class WebPageContentSplitter:
    content_type = "text/html"

    def __init__(
        self, content: WebPageContent, *, chunk_size: int = 1024, chunk_overlap=0
    ):
        self.content = content
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap

        self.chunks = []

    def _build_metadata(self, metadata: dict | None, **kwargs) -> dict:
        """Remove None values from metadata"""
        cleaned = {k: v for k, v in kwargs.items() if v is not None}
        if metadata is None:
            d = dict(**cleaned)
            return d
        d = {k: v for k, v in metadata.items() if v is not None}
        d.update(cleaned)
        return d

    def _split(self) -> List[Document]:
        splitter = RecursiveCharacterTextSplitter(
            chunk_size=self.chunk_size, chunk_overlap=self.chunk_overlap
        )
        metadata = self._build_metadata(
            self.content.metadata,
        )
        docs = [Document(page_content=self.content.content, metadata=metadata)]
        splits = splitter.split_documents(documents=docs)
        return splits

    def split(self) -> List[ChunkNode]:
        splits = self._split()
        if len(self.chunks) > 0:
            logger.warning("existing chunks will be cleared.")
            self.chunks = []
        prev = None
        for i, item in enumerate(splits):
            chunk = ChunkNode(
                uid=self.content.uid,
                source=self.content.source,
                uri=self.content.uri,
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
        chunks = [chunk.to_dict() for chunk in self.chunks]
        return {
            "uid": self.content.uid,
            "source": self.content.source,
            "uri": self.content.uri,
            "chunk_size": self.chunk_size,
            "chunk_overlap": self.chunk_overlap,
            "chunks": chunks,
        }
