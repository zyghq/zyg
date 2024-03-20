from uuid import uuid4

from langchain_community.embeddings import OllamaEmbeddings
from langchain_community.vectorstores import Chroma
from langchain_core.documents import Document

from indexer.loaders.parsers import WebPageParser
from indexer.splitters.splitter import WebDocumentSplitter

VECTOR_DB = "./vectordb"


def clean_metadata(metadata: dict) -> dict:
    return {k: v for k, v in metadata.items() if v is not None}


if __name__ == "__main__":
    local_file_path = "./output.html"
    uid = str(uuid4())
    parser = WebPageParser(uid, local_file_path)
    content = parser.parse()

    splitter = WebDocumentSplitter(content)
    chunks = splitter.split()

    docs = [
        Document(
            page_content=chunk.content,
            metadata={
                **chunk.metadata,
                "counter": chunk.chunk_counter,
                **clean_metadata(content.metadata),
            },
        )
        for chunk in chunks
    ]
    ids = [chunk.chunk_id for chunk in chunks]
    vectorstore = Chroma.from_documents(
        documents=docs,
        embedding=OllamaEmbeddings(model="nomic-embed-text"),
        ids=ids,
        collection_name="devcollection",
        persist_directory=VECTOR_DB,
    )
