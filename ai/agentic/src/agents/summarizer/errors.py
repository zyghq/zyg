class PromptError(Exception):
    """Raised when there is an error with prompt"""

    def __init__(self, detail: str):
        super().__init__(f"prompt error detail={detail}")
        self.detail = detail


class SummarizerError(Exception):
    """Raised in general when there is error when generating summary"""

    def __init__(self, detail: str) -> None:
        super().__init__(f"summary error detail={detail}")
        self.detail = detail
