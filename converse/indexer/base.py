from typing import Dict

from pydantic import BaseModel


class WebPageContent(BaseModel):
    """Represents required attributes of a page."""

    uid: str  # uniquely identifies the web content.
    uri: str  # internal uri of the web content.
    source: str  # actual source of hosted content. e.g. https://example.com/blog/

    content: str
    content_type: str = "text/html"  # represents the web content type
    metadata: Dict[str, str | int] = {}  # additional information about the content


class ChunkNode:
    def __init__(
        self,
        *,
        uid: str,
        source: str,
        uri: str,
        chunk_id: str,
        chunk_counter: int,
        content: str,
        metadata: Dict[str, str | int],
        parent: "ChunkNode | None" = None,
        child: "ChunkNode | None" = None,
    ):
        self.uid = uid
        self.source = source
        self.uri = uri
        self.chunk_id = chunk_id
        self.chunk_counter = chunk_counter
        self.content = content
        self.metadata = metadata
        self.parent = parent
        self.child = child

        assert uid is not None, "uid is required"
        assert source is not None, "source is required"
        assert uri is not None, "uri is required"
        assert chunk_id is not None, "chunk_id is required"
        assert chunk_counter is not None, "chunk_counter is required"
        assert content is not None, "content is required"

    def to_dict(self) -> dict:
        return {
            "uid": self.uid,
            "source": self.source,
            "uri": self.uri,
            "chunk_id": self.chunk_id,
            "chunk_counter": self.chunk_counter,
            "content": self.content,
            "metadata": self.metadata,
            "parent": self.parent.chunk_id if self.parent else None,
            "child": self.child.chunk_id if self.child else None,
        }
