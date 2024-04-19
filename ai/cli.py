import typer
from langchain import hub
from langchain_community.embeddings import OllamaEmbeddings
from langchain_community.llms import Ollama
from langchain_community.vectorstores import Chroma
from langchain_core.output_parsers import StrOutputParser
from langchain_core.runnables import RunnablePassthrough
from rich.progress import Progress, SpinnerColumn, TextColumn

from command import WebPageIndexerCommand

app = typer.Typer()


_CHROMA_VECTOR_DIR = "./vectordb"
_CHROMA_DEFAULT_COLLECTION = "devcollection"


vectorstore = Chroma(
    collection_name=_CHROMA_DEFAULT_COLLECTION,
    embedding_function=OllamaEmbeddings(model="nomic-embed-text"),
    persist_directory=_CHROMA_VECTOR_DIR,
)


@app.command()
def main():
    typer.echo("Hello, World!")


@app.command()
def index(
    type: str = typer.Argument("type", help="Type of the content to index"),
    source: str = typer.Option("--source ", help="Source of the content to index"),
    uri: str = typer.Option("--uri", help="Internal URI of the content"),
    model: str = typer.Option("--model", help="Model to use for embedding"),
):
    if type == "web":
        command = WebPageIndexerCommand(source=source, uri=uri)
        command.run(model, save=True)
    else:
        typer.echo("unsupported index type")


@app.command()
def ama():
    query = typer.prompt("", prompt_suffix=">> ")
    retriever = vectorstore.as_retriever(search_kwargs={"k": 3})
    prompt = hub.pull("rlm/rag-prompt")
    llm = Ollama(model="llama2")

    chain = (
        {"context": retriever, "question": RunnablePassthrough()}
        | prompt
        | llm
        | StrOutputParser()
    )

    with Progress(
        SpinnerColumn(),
        TextColumn("[progress.description]{task.description}"),
        transient=True,
    ) as progress:
        progress.add_task("Thinking...", total=None)
        response = chain.invoke(query)
    typer.echo(response)


if __name__ == "__main__":
    app()
