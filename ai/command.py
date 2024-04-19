import json
import logging
import time
from concurrent.futures import ThreadPoolExecutor
from typing import Dict
from uuid import uuid4

import chromadb
from langchain import hub
from langchain_community.embeddings import OllamaEmbeddings
from langchain_community.llms import Ollama
from langchain_community.vectorstores import Chroma
from langchain_core.output_parsers import StrOutputParser
from langchain_core.runnables import RunnablePassthrough

from indexer.embedders.embedder import LocalOllamaEmbedding
from indexer.loaders.parsers import WebPageParser
from indexer.splitters.splitter import WebPageContentSplitter

logger = logging.getLogger(__name__)


_CHROMA_VECTOR_DIR = "./vectordb"
_CHROMA_DEFAULT_COLLECTION = "devcollection"


vectorstore = Chroma(
    collection_name=_CHROMA_DEFAULT_COLLECTION,
    embedding_function=OllamaEmbeddings(model="nomic-embed-text"),
    persist_directory=_CHROMA_VECTOR_DIR,
)


client = chromadb.PersistentClient(_CHROMA_VECTOR_DIR)


def build_metadata(metadata: Dict[str, str] | None, **kwargs) -> Dict[str, str]:
    """Remove None values from metadata"""
    cleaned = {k: v for k, v in kwargs.items() if v is not None}
    if metadata is None:
        d = dict(**cleaned)
        return d
    d = {k: v for k, v in metadata.items() if v is not None}
    d.update(cleaned)
    return d


class WebPageIndexerCommand:

    def __init__(self, *, source: str, uri: str, uid=str(uuid4())) -> None:
        self.source = source
        self.uri = uri
        self.uid = uid

    def run(self, embed_model: str, save: bool = False):
        start = time.time()

        parser = WebPageParser(self.uid, source=self.source, uri=self.uri)
        content = parser.parse()

        splitter = WebPageContentSplitter(content)
        chunks = splitter.split()

        embedder = LocalOllamaEmbedding(embed_model)

        if save and len(chunks) > 0:
            with open("./store/splits.json", "w") as f:
                json.dump(splitter.to_dict(), f, indent=2)
                logger.info("saved splits to output.json")

        def embeds(chunk):
            result = embedder(chunk)
            return result

        with ThreadPoolExecutor(max_workers=4) as executor:
            embeddings = list(executor.map(embeds, chunks))

        if save:
            with open("./store/embeddings.json", "w") as f:
                json.dump([e.to_dict() for e in embeddings], f, indent=2)
                logger.info("saved chunk embeddings to embeddings.json")

        collection = client.get_or_create_collection(
            name="devcollection", metadata={"hnsw:space": "cosine"}
        )

        ids = [e.chunk_id for e in embeddings]
        documents = [e.content for e in chunks]
        metadatas = [
            build_metadata(
                None,
                document_uid=e.document_uid,
                document_url=e.document_url,
                chunk_id=e.chunk_id,
                chunk_counter=e.chunk_counter,
                parent=e.parent,
                child=e.child,
            )
            for e in embeddings
        ]
        chroma_embeddings = [e.embedding for e in embeddings]

        collection.add(
            documents=documents,
            embeddings=chroma_embeddings,
            ids=ids,
            metadatas=list(metadatas),
        )

        logger.info("eta: ", time.time() - start)


class QueryCommand:

    def __init__(self) -> None:
        pass

    def query(self, q: str) -> str:
        retriever = vectorstore.as_retriever(search_kwargs={"k": 3})
        prompt = hub.pull("rlm/rag-prompt")
        llm = Ollama(model="llama2")

        chain = (
            {"context": retriever, "question": RunnablePassthrough()}
            | prompt
            | llm
            | StrOutputParser()
        )

        response = chain.invoke(q)
        return response
