from typing import Any

import requests
from bs4 import BeautifulSoup

from indexer.base import WebPageContent


class WebPageParser:
    def __init__(self, uid: str, *, source: str, uri: str) -> None:
        self.uid = uid
        self.source = source
        self.uri = uri

    def _feth(self) -> str:
        """fetch web page content"""
        response = requests.get(self.source)
        if response.status_code >= 200 and response.status_code < 300:
            return response.text
        raise ValueError(
            f"failed to fetch content from {self.source} with internal uri {self.uri}"
        )

    @staticmethod
    def _build_metadata(soup: Any, url: str) -> dict:
        """Build metadata from BeautifulSoup output for the HTML page."""
        metadata = {"source": url}
        if title := soup.find("title"):
            metadata["title"] = title.get_text()
        if description := soup.find("meta", attrs={"name": "description"}):
            metadata["description"] = description.get(
                "content", "No description found."
            )
        if html := soup.find("html"):
            metadata["language"] = html.get("lang", "No language found.")
        return metadata

    def parse(self) -> WebPageContent:
        """entry point for parsing web page"""
        html = self._feth()
        soup = BeautifulSoup(html, "html.parser")
        metadata = self._build_metadata(soup, self.source)
        for match in soup(["script", "style", "a"]):
            match.decompose()
        texts = [element.get_text(separator="\n", strip=True) for element in soup]
        content = "\n".join(texts)
        return WebPageContent(
            uid=self.uid,
            uri=self.uri,
            source=self.source,
            content=content,
            metadata=metadata,
        )
