import base64
import os
import webbrowser

import inquirer
import pkg_resources
import typer

from ..alerting.handlers import alert_drift_handler
from ..alerting.transport import AbstractAlertTransport, ConsoleAlertTransport
from ..dbt.snapshot import (get_snapshot_dates, get_snapshot_diff,
                            get_snapshot_nodes)
from ..dbt.snapshot_to_drift import convert_snapshot_to_drift_summary
from ..logger import get_logger
from ..user_defined_function import import_user_defined_function
from .common import get_user_date_selection

app = typer.Typer()

logger = get_logger(__name__)


@app.command()
def show(snapshot_id: str = typer.Option(None, help="id of your snapshot")):
    snapshot_nodes = get_snapshot_nodes()
    snapshot_node = get_or_prompt_snapshot_node(snapshot_id, snapshot_nodes)
    snapshot_dates = get_snapshot_dates(snapshot_node)

    snapshot_date = get_user_date_selection(snapshot_dates)
    if snapshot_date is None:
        typer.echo("No snapshot data for selected date. Exiting.")
        raise typer.Exit(code=1)

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


@app.command()
def check(snapshot_id: str = typer.Option(None, help="id of your snapshot"), date: str = typer.Option(None, help="date of your snapshot")):
    snapshot_node = get_or_prompt_snapshot_node(snapshot_id, get_snapshot_nodes())
    snapshot_date = get_user_date_selection(get_snapshot_dates(snapshot_node), date)

    if snapshot_date is None:
        typer.echo("No snapshot data for selected date. Exiting.")
        raise typer.Exit(code=1)

    print(f"Getting {snapshot_node['unique_id']} for {snapshot_date}.")

    [drift_handler, alert_transport] = get_user_defined_handlers(snapshot_node)

    if not isinstance(alert_transport, AbstractAlertTransport):
        print("Alert transport is not an instance of AbstractAlertTransport, defaulting to ConsoleAlertTransport.")
        alert_transport = ConsoleAlertTransport()


    diff = get_snapshot_diff(snapshot_node, snapshot_date)
    context = convert_snapshot_to_drift_summary(snapshot_diff=diff, id_column="month", date_column="month")
    alert = drift_handler(context)
    alert_title = f"Drift alert for {snapshot_node['unique_id']} on {snapshot_date}"
    alert_transport.send(alert_title, alert, context)



def get_user_defined_handlers(snapshot_node):
    snapshot_file_path = snapshot_node["original_file_path"]
    directory_path = os.path.dirname(snapshot_file_path)
    snapshot_file_name = os.path.basename(snapshot_file_path)
    snapshot_file_name_without_extension, _ = os.path.splitext(snapshot_file_name)
    user_defined_file_path = os.path.join(directory_path, f"{snapshot_file_name_without_extension}.datadrift.py")
    
    print(f"Looking for user defined handlers in {user_defined_file_path}")

    try:
        [drift_handler, alert_transport] = import_user_defined_function(user_defined_file_path,[ "drift_handler", "alert_transport"])
        return [drift_handler, alert_transport]
    except Exception as e:
        logger.error(f"Error importing user defined handler: {e}")
        logger.warn("No user defined handler found. Using default handler.")
        return [alert_drift_handler, None]


def get_or_prompt_snapshot_node(snapshot_id, snapshot_nodes):
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
    return snapshot_node


if __name__ == "__main__":
    app()
