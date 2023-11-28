from datetime import date, timedelta

import inquirer
import typer
from driftdb.dbt.snapshot import get_snapshot_dates, get_snapshot_nodes

from .common import get_user_date_selection

app = typer.Typer()


@app.command()
def show(snapshot_id: str = typer.Option(None, help="id of your snapshot")):
    snapshot_nodes = get_snapshot_nodes()
    if not snapshot_id:
        questions = [
            inquirer.List(
                "choice",
                message="Please choose a snapshot to show",
                choices=[node["unique_id"] for node in snapshot_nodes],
            ),
        ]
        answers = inquirer.prompt(questions)
        if answers is None:
            typer.echo("No choice selected. Exiting.")
            raise typer.Exit(code=1)

        snapshot_id = answers["choice"]

    snapshot_node = [node for node in snapshot_nodes if node["unique_id"] == snapshot_id][0]
    snapshot_dates = get_snapshot_dates(snapshot_node)

    snapshot_date = get_user_date_selection(snapshot_dates)
    print(snapshot_date)
    print(type(snapshot_date))
