from typing import Any

from bs4 import BeautifulSoup

from core.base import WebPageContent


class WebPageParser:
    def __init__(self, uid, url) -> None:
        self.uid = uid
        self.url = url

    def _feth(self) -> str:
        """fetch web page content"""
        # TODO: update this method to fetch content from the web
        with open(self.url, "r", encoding="utf-8") as f:
            return f.read()

    @staticmethod
    def _build_metadata(soup: Any, url: str) -> dict:
        """Build metadata from BeautifulSoup output."""
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
        metadata = self._build_metadata(soup, self.url)
        for match in soup(["script", "style", "a"]):
            match.decompose()
        texts = [element.get_text(separator="\n", strip=True) for element in soup]
        content = "\n".join(texts)
        return WebPageContent(
            uid=self.uid, url=self.url, metadata=metadata, content=content
        )
