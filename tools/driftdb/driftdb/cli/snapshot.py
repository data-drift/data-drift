import inquirer
import typer
from driftdb.dbt.snapshot import get_snapshot_nodes

app = typer.Typer()


@app.command()
def show(name: str = typer.Option(None, help="name of your snapshot")):
    if not name:
        snapshot_nodes = get_snapshot_nodes()
        questions = [
            inquirer.List(
                "choice",
                message="Please choose a snapshot to show",
                choices=snapshot_nodes,
            ),
        ]
        answers = inquirer.prompt(questions)
        if answers is None:
            typer.echo("No choice selected. Exiting.")
            raise typer.Exit(code=1)

        name = answers["choice"]

    typer.echo(name)
