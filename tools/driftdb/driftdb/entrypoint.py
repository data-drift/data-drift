import typer
from driftdb.cli.server import start_server

from . import version
from .cli.dbt import app as dbt
from .cli.seed import app as seed
from .cli.snapshot import app as snapshot
from .cli.store import app as store

app = typer.Typer()

app.add_typer(dbt, name="dbt")
app.add_typer(seed, name="seed")
app.add_typer(store, name="store")
app.add_typer(snapshot, name="snapshot")


def print_version(value: bool):
    if value:
        typer.echo(f"Driftdb CLI Version: {version.version}")  # Use the imported version variable
        raise typer.Exit()


# Define a callback function for the app
@app.callback()
def main(
    version: bool = typer.Option(
        None,
        "--version",
        "-v",
        help="Show the version",
        is_eager=True,  # Ensures that this option is handled first
        callback=print_version,
    )
):
    """
    Your CLI's main entry point.
    """
    if version:
        raise typer.Exit()


@app.command()
def start():
    start_server()


if __name__ == "__main__":
    app()
