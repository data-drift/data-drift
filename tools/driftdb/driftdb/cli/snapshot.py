import base64
import os
import webbrowser

import inquirer
import pkg_resources
import typer
from driftdb.dbt.snapshot import get_snapshot_dates, get_snapshot_diff, get_snapshot_nodes

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

    diff = get_snapshot_diff(snapshot_node, snapshot_date)

    spa_html_path = pkg_resources.resource_filename(__name__, "../spa/snapshot/index.html")

    with open(spa_html_path, "r", encoding="utf-8") as spa_html_file:
        spa_html_code = spa_html_file.read()

    json_diff = diff.to_json(date_format="iso")

    encoded_diff = base64.b64encode(json_diff.encode("utf-8"))
    b64string_diff = encoded_diff.decode("utf-8")
    compiled_output_html = (
        f"<script>" f"window.generated_diff = JSON.parse(atob('{b64string_diff}'));" f"</script>" f"{spa_html_code}"
    )

    with open("diff.html", "w") as f:
        f.write(compiled_output_html)

    html_file_path = os.path.abspath("diff.html")

    webbrowser.open("file://" + html_file_path)


if __name__ == "__main__":
    app()
