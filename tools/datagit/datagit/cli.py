import click
import json
import pandas as pd
from datagit import github_connector
from datagit.drift_evaluators import auto_merge_drift
from github import Github


@click.group()
def cli_entrypoint():
    pass


@cli_entrypoint.group()
def dbt():
    pass


@dbt.command()
@click.option("--token", prompt="Your token", help="Number of greetings.")
@click.option("--repo", prompt="Your repo", help="The person to greet.")
def run(token, repo):
    import dbt
    from dbt.cli.main import dbtRunner
    from dbt.config.runtime import load_profile, load_project, RuntimeConfig
    import dbt.version
    from dbt.adapters.factory import get_adapter

    """Simple program that greets NAME for a total of COUNT times."""
    click.echo(f"Repo is {repo}!")

    project_path = "."
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
            table.to_csv("data.csv")

            dataframe = pd.read_csv("data.csv")
            github_connector.store_metric(
                dataframe=dataframe,
                ghClient=Github(token),
                branch="main",
                filepath=repo + "/dbt-drift/metrics/" + node["name"] + ".csv",
                drift_evaluator=auto_merge_drift,
            )
