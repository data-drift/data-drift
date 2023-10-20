import sys
import threading
import webbrowser
import click
import json
from datagit.dataset import generate_dataframe, insert_drift
import pandas as pd
import os
from datagit import github_connector
from datagit import local_connector
from datagit.drift_evaluators import auto_merge_drift
from github import Github
import subprocess
import platform

import pkg_resources
import http.server
import socketserver


@click.group()
def cli_entrypoint():
    pass


@cli_entrypoint.group()
def dbt():
    pass


@dbt.command()
@click.option(
    "--token",
    default=lambda: os.environ.get("DATADRIFT_GITHUB_TOKEN", ""),
    help="Token to access your repo. With PR and Content read and write rights",
)
@click.option(
    "--repo",
    default=lambda: os.environ.get("DATADRIFT_GITHUB_REPO", ""),
    help="The datagit repo in the form org/repo",
)
@click.option("--project-dir", default=".", help="The dbt project dir")
def run(token, repo, project_dir):
    from dbt.cli.main import dbtRunner
    from dbt.config.runtime import load_profile, load_project, RuntimeConfig
    from dbt.adapters.factory import get_adapter

    if not repo:
        repo = click.prompt("Your repo")

    if not token:
        token = click.prompt("Your token")
    click.echo(f"Pushing to {repo}!")

    project_path = project_dir
    dbtRunner().invoke(["-q", "debug"], project_dir=str(project_path))
    profile = load_profile(str(project_path), {})
    project = load_project(str(project_path), version_check=False, profile=profile)

    runtime_config = RuntimeConfig.from_parts(project, profile, {})

    adapter = get_adapter(runtime_config)

    with open(f"{project_path}/target/manifest.json") as manifest_file:
        manifest = json.load(manifest_file)

    data_drift_nodes = [
        node
        for node in manifest["nodes"].values()
        if node["config"]["meta"]["datadrift"]
    ]

    for node in data_drift_nodes:
        query = f'SELECT {node["config"]["meta"]["datadrift_unique_key"]} as unique_key,{node["config"]["meta"]["datadrift_date"]} as date, * FROM {node["relation_name"]}'
        with adapter.connection_named("default"):
            resp, table = adapter.execute(query, fetch=True)

            #  TODO: Try to transform table in a dataframe without writing to a file
            metric_file = "data.csv"
            table.to_csv(metric_file)
            dataframe = pd.read_csv(metric_file)

            github_connector.store_metric(
                dataframe=dataframe,
                ghClient=Github(token),
                branch="main",
                filepath=repo + "/dbt-drift/metrics/" + node["name"] + ".csv",
                drift_evaluator=auto_merge_drift,
            )

            os.remove(metric_file)


@cli_entrypoint.command()
def start():
    click.echo("Starting the application...")

    if platform.system() == "Darwin":
        if platform.machine().startswith("arm"):
            binary_path = pkg_resources.resource_filename(
                "datagit", "bin/data-drift-mac-m1"
            )
        else:
            binary_path = pkg_resources.resource_filename(
                "datagit", "bin/data-drift-mac-intel"
            )
    else:
        # TODO: Update this path for other platforms (Linux, Windows, etc.)
        raise ValueError("Unsupported platform")

    # Get a copy of the current environment variables
    env = os.environ.copy()

    # Set the PORT environment variable
    env["PORT"] = "9740"

    server_process = subprocess.Popen(
        [binary_path],
        env=env,
    )

    PORT = 9741

    SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
    DIRECTORY = os.path.join(SCRIPT_DIR, "bin/frontend/dist")

    class Handler(http.server.SimpleHTTPRequestHandler):
        def __init__(self, *args, **kwargs):
            super().__init__(*args, directory=DIRECTORY, **kwargs)

        def do_GET(self):
            print(f"Request path: {self.path}")
            # If the requested URL maps to an existing file, serve that.
            if os.path.exists(self.translate_path(self.path)):
                super().do_GET()
                return

            # Otherwise, serve the main index.html file.
            self.path = "index.html"
            super().do_GET()

    httpd = socketserver.TCPServer(("", PORT), Handler)

    try:
        print(f"Serving directory '{DIRECTORY}' on port {PORT}")
        url = f"http://localhost:{PORT}/tables"
        print("Opening browser...", url)
        webbrowser.open(url)
        httpd.serve_forever()
        server_process.wait()

    except KeyboardInterrupt:
        click.echo("Shutting down servers...")
        httpd.shutdown()
        click.echo("Httpd shut down")
        server_process.terminate()
        click.echo("Server down")
        sys.exit()


@cli_entrypoint.group()
def seed():
    pass


@seed.command()
@click.option(
    "--table",
    help="name of your table",
)
@click.option(
    "--row-number",
    default=10000,
    help="Number of line to generate",
)
def create(table, row_number):
    if not table:
        table = click.prompt("Table name")
    click.echo("Creating seed file...")
    dataframe = generate_dataframe(row_number)
    click.echo(dataframe)
    local_connector.store_metric(metric_name=table, metric_value=dataframe)
    click.echo("Creating seed created...")


@seed.command()
@click.option(
    "--table",
    help="name of your table",
)
@click.option(
    "--row-number",
    default=100,
    help="Number of line to update",
)
def update(table, row_number):
    if not table:
        tables = local_connector.get_metrics()
        table = select_from_list("Please enter table number", tables)

    click.echo("Updating seed file...")
    dataframe = local_connector.get_metric(metric_name=table)
    drifted_dataset = insert_drift(dataframe, row_number)
    local_connector.store_metric(metric_name=table, metric_value=drifted_dataset)


@cli_entrypoint.command()
@click.argument("csvpathfile")
@click.option(
    "--table",
    help="name of your table",
)
@click.option(
    "--unique-key-column",
    help="name of your unique key column",
)
@click.option(
    "--date-column",
    help="name of your date column",
)
def load_csv(csvpathfile, table, unique_key_column, date_column):
    if not table:
        tables = local_connector.get_metrics()
        table = click.prompt(
            "Please enter table name (exising or not)", type=click.Choice(tables)
        )

    click.echo(f"Loading CSV file {csvpathfile}...")
    assert os.path.exists(csvpathfile), f"CSV file {csvpathfile} does not exist"
    dataframe = pd.read_csv(csvpathfile)

    if "unique_key" not in dataframe.columns:
        if not unique_key_column:
            unique_key_column = click.prompt(
                "Please enter unique key column name",
                type=click.Choice(dataframe.columns),
            )

        assert (
            unique_key_column in dataframe.columns
        ), f"Column {unique_key_column} does not exist in CSV file"
        dataframe.insert(0, "unique_key", dataframe[unique_key_column])

    if "date" not in dataframe.columns:
        if not date_column:
            date_column = click.prompt(
                "Please enter date column name", type=click.Choice(dataframe.columns)
            )
        assert (
            date_column in dataframe.columns
        ), f"Column {date_column} does not exist in CSV file"
        dataframe.insert(1, "date", dataframe[date_column])

    local_connector.store_metric(metric_name=table, metric_value=dataframe)


def select_from_list(prompt, choices):
    for idx, item in enumerate(choices, 1):
        click.echo(f"{idx}: {item}")
    selection = click.prompt(prompt, type=click.IntRange(1, len(choices)))
    return choices[selection - 1]
