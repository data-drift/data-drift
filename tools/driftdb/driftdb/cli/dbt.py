import json
import os
from datetime import datetime

import click
import numpy as np
import pandas as pd
import pytz
import typer
from driftdb.cli.server import start_server
from driftdb.connectors.github_connector import GithubConnector
from driftdb.connectors.local_connector import LocalConnector
from github.MainClass import Github

from ..logger import get_logger
from .common import dbt_adapter_query, prompt_from_list

app = typer.Typer()

logger = get_logger(__name__)


@app.command()
def run(
    token: str = typer.Option(
        os.environ.get("DATADRIFT_GITHUB_TOKEN", ""),
        help="Token to access your repo. With PR and Content read and write rights",
    ),
    repo: str = typer.Option(os.environ.get("DATADRIFT_GITHUB_REPO", ""), help="The driftdb repo in the form org/repo"),
    storage: str = typer.Option("local", help="Whether to use local or github storage"),
    project_dir: str = typer.Option(".", help="The dbt project dir"),
):
    from dbt.adapters.factory import get_adapter
    from dbt.cli.main import dbtRunner
    from dbt.config.runtime import RuntimeConfig, load_profile, load_project

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

    click.echo(f"Parsing manifest")
    with open(f"{project_path}/target/manifest.json") as manifest_file:
        manifest = json.load(manifest_file)

    data_drift_nodes = [
        node
        for node in manifest["nodes"].values()
        if (node["resource_type"] == "model") & node["config"]["meta"].get("datadrift", False)
    ]

    for node in data_drift_nodes:
        query = f'SELECT {node["config"]["meta"]["datadrift_unique_key"]} as unique_key,{node["config"]["meta"]["datadrift_date"]} as date, * FROM {node["relation_name"]}'
        with adapter.connection_named("default"):  # type: ignore
            dataframe = dbt_adapter_query(adapter, query)

            if storage == "github":
                github_connector = GithubConnector(
                    github_client=Github(token),
                    github_repository_name=repo,
                )
                github_connector.snapshot_table(
                    table_dataframe=dataframe,
                    table_name="/dbt-drift/metrics/" + node["name"] + ".csv",
                )
            else:
                local_connector = LocalConnector()
                local_connector.snapshot_table(
                    table_name=node["name"],
                    table_dataframe=dataframe,
                )


@app.command()
def snapshot():
    from dbt.adapters.factory import get_adapter
    from dbt.cli.main import dbtRunner
    from dbt.config.runtime import RuntimeConfig, load_profile, load_project

    project_dir = "."
    project_path = project_dir
    dbtRunner().invoke(["-q", "debug"], project_dir=str(project_path))
    profile = load_profile(str(project_path), {})
    project = load_project(str(project_path), version_check=False, profile=profile)

    runtime_config = RuntimeConfig.from_parts(project, profile, {})

    adapter = get_adapter(runtime_config)

    with open(f"{project_path}/target/manifest.json") as manifest_file:
        manifest = json.load(manifest_file)

    snapshot_nodes = [node for node in manifest["nodes"].values() if (node["resource_type"] == "snapshot")]

    [snapshot_name, snapshot_index] = prompt_from_list(
        "Which snapshot do you want to process?",
        [node["name"] for node in snapshot_nodes],
    )

    node = snapshot_nodes[snapshot_index]
    logger.info(f"Handling node: {node['unique_id']}")
    snapshot_table = node["relation_name"]
    default_date_column = "created_at"

    # prompt the user for a date column
    date_column = click.prompt("Enter the date column name", default=default_date_column)

    unique_key = node["config"]["unique_key"]
    metric_name = node["name"]

    with adapter.connection_named("default"):  # type: ignore
        text_query = f"""
        SELECT DISTINCT dbt_valid_from FROM {snapshot_table}
        UNION
        SELECT DISTINCT dbt_valid_to FROM {snapshot_table} WHERE NOT NULL;
        """

        df = dbt_adapter_query(adapter, text_query)
        df["dbt_valid_from"] = pd.to_datetime(df["dbt_valid_from"])

        local_connector = LocalConnector()
        metric_history = local_connector.get_table_history(table_name=metric_name)
        latest_commit = next(metric_history, None)
        if latest_commit is not None:
            authored_date = latest_commit.authored_date
            authored_datetime = datetime.fromtimestamp(authored_date)
            # Compare only the second part, without the milliseconds, and take dates after the latest measure commit
            df = df.loc[df["dbt_valid_from"].dt.floor("S") > authored_datetime]

        logger.info(f"Snapshot dates to process: {df}")

        for index, row in df.iterrows():
            date = row["dbt_valid_from"]
            logger.info(f"Processing data for date: {date}")
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
            logger.info(data_as_of_date)

            # Compute date in UTC format
            localized_date = date.replace(tzinfo=pytz.UTC)

            local_connector = LocalConnector()
            local_connector.snapshot_table(
                table_name=metric_name,
                table_dataframe=data_as_of_date,
                measure_date=localized_date,
            )
    start_server("/tables/" + metric_name)
