import sys
import webbrowser
import click
import json
from datagit.dataset import generate_dataframe, insert_drift
import numpy as np
import pandas as pd
import os
from datagit import github_connector
from datagit import local_connector
from datagit.drift_evaluators import auto_merge_drift
from github import Github
import subprocess
import platform
from datetime import datetime
import pkg_resources
import http.server
import socketserver
from tzlocal import get_localzone


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
@click.option(
    "--storage",
    default="local",
    help="Wether to use local or github storage",
)
@click.option("--project-dir", default=".", help="The dbt project dir")
def run(token, repo, storage, project_dir):
    from dbt.cli.main import dbtRunner
    from dbt.config.runtime import load_profile, load_project, RuntimeConfig
    from dbt.adapters.factory import get_adapter

    if storage == "github":
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
        if (node["resource_type"] == "model")
        & node["config"]["meta"].get("datadrift", False)
    ]

    for node in data_drift_nodes:
        query = f'SELECT {node["config"]["meta"]["datadrift_unique_key"]} as unique_key,{node["config"]["meta"]["datadrift_date"]} as date, * FROM {node["relation_name"]}'
        with adapter.connection_named("default"):
            dataframe = dbt_adapter_query(adapter, query)

            if storage == "github":
                github_connector.store_metric(
                    dataframe=dataframe,
                    ghClient=Github(token),
                    branch="main",
                    filepath=repo + "/dbt-drift/metrics/" + node["name"] + ".csv",
                    drift_evaluator=auto_merge_drift,
                )
            else:
                local_connector.store_metric(
                    metric_name=node["name"], metric_value=dataframe
                )


@dbt.command()
def snapshot():
    from dbt.cli.main import dbtRunner
    from dbt.config.runtime import load_profile, load_project, RuntimeConfig
    from dbt.adapters.factory import get_adapter

    project_dir = "."
    project_path = project_dir
    dbtRunner().invoke(["-q", "debug"], project_dir=str(project_path))
    profile = load_profile(str(project_path), {})
    project = load_project(str(project_path), version_check=False, profile=profile)

    runtime_config = RuntimeConfig.from_parts(project, profile, {})

    adapter = get_adapter(runtime_config)

    with open(f"{project_path}/target/manifest.json") as manifest_file:
        manifest = json.load(manifest_file)

    snapshot_nodes = [
        node
        for node in manifest["nodes"].values()
        if (node["resource_type"] == "snapshot")
    ]

    for node in snapshot_nodes:
        print("Handling node:", node["unique_id"])
        with adapter.connection_named("default"):
            snapshot_table = node["relation_name"]
            date_column = node["config"]["meta"]["datadrift_date"]
            unique_key = node["config"]["unique_key"]
            metric_name = node["name"]

            text_query = f"""
            SELECT DISTINCT dbt_valid_from FROM {snapshot_table}
            UNION
            SELECT DISTINCT dbt_valid_to FROM {snapshot_table} WHERE NOT NULL;
            """

            df = dbt_adapter_query(adapter, text_query)
            df["dbt_valid_from"] = pd.to_datetime(df["dbt_valid_from"])

            metric_history = local_connector.get_metric_history(metric_name=metric_name)
            latest_commit = next(metric_history, None)
            if latest_commit is not None:
                authored_date = latest_commit.authored_date
                authored_datetime = datetime.fromtimestamp(authored_date)
                # Compare only the second part, without the milliseconds, and take dates after the latest measure commit
                df = df.loc[df["dbt_valid_from"].dt.floor("S") > authored_datetime]

            print("Snapshot dates to process:", df)

            for index, row in df.iterrows():
                date = row["dbt_valid_from"]
                print(f"Processing data for date: {date}")
                if pd.isna(date):
                    date = pd.Timestamp.now()

                asOfQuery = f"""
                SELECT * FROM {snapshot_table}
                WHERE ('{date}' >= dbt_valid_from AND  '{date}' < dbt_valid_to) OR
                    (dbt_valid_to IS NULL AND dbt_valid_from <= '{date}');
                """

                data_as_of_date = dbt_adapter_query(adapter, asOfQuery)

                data_as_of_date.replace({np.nan: "NA"}, inplace=True)

                # Drop dbt columns and add date and unique_key columns
                data_as_of_date = data_as_of_date.drop(
                    ["dbt_scd_id", "dbt_updated_at", "dbt_valid_from", "dbt_valid_to"],
                    axis=1,
                )
                data_as_of_date["date"] = data_as_of_date[date_column]
                data_as_of_date["unique_key"] = data_as_of_date[unique_key]
                print(data_as_of_date)

                # Compute date in UTC format
                local_tz = get_localzone()
                localized_date = date.replace(tzinfo=local_tz)

                local_connector.store_metric(
                    metric_name=metric_name,
                    metric_value=data_as_of_date,
                    measure_date=localized_date,
                )


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
@click.option(
    "--table",
    help="name of your table",
)
def delete(table):
    tables = local_connector.get_metrics()
    if not table:
        tables = local_connector.get_metrics()
        table = select_from_list("Please enter table number", tables)
    local_connector.delete_metric(metric_name=table)


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


def dbt_adapter_query(
    adapter,
    query: str,
) -> pd.DataFrame:
    _, table = adapter.execute(query, fetch=True)
    data = {
        column_name: table.columns[column_name].values()
        for column_name in table.column_names
    }
    return pd.DataFrame(data)
