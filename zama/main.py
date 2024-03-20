from uuid import uuid4

from indexers.indexer import WebDocumentSplitter
from loaders.parsers import WebPageParser

if __name__ == "__main__":
    local_file_path = "./output.html"
    uid = str(uuid4())
    parser = WebPageParser(uid, local_file_path)
    content = parser.parse()

    splitter = WebDocumentSplitter(content)
    splitter.split()
    print(splitter.to_dict())
