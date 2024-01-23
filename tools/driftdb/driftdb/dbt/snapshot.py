import json

from pandas import DataFrame
from typing_extensions import List, TypedDict

from ..cli.common import dbt_adapter_query


class SnapshotNode(TypedDict):
    unique_id: str
    name: str
    relation_name: str


def get_snapshot_nodes() -> List[SnapshotNode]:
    project_dir = "."

    with open(f"{project_dir}/target/manifest.json") as manifest_file:
        manifest = json.load(manifest_file)

        snapshot_nodes = [node for node in manifest["nodes"].values() if (node["resource_type"] == "snapshot")]
        return snapshot_nodes


def get_snapshot_dates(snapshot_node: SnapshotNode) -> List[str]:
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

    with adapter.connection_named("default"):  # type: ignore
        text_query = f"""
        SELECT DISTINCT dbt_valid_from
        FROM {snapshot_node["relation_name"]}
        WHERE dbt_valid_from IS NOT NULL
        ORDER BY dbt_valid_from DESC;
        """

        df = dbt_adapter_query(adapter, text_query)
        date_strings = df["dbt_valid_from"].dt.strftime("%Y-%m-%d %H:%M:%S.%f").tolist()
        return date_strings


def get_snapshot_diff(snapshot_node: SnapshotNode, snapshot_date: str) -> DataFrame:
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

    with adapter.connection_named("default"):  # type: ignore
        text_query = f"""
        WITH 
        valid_from AS(
            SELECT *,
                'after' AS record_status
            FROM {snapshot_node["relation_name"]}
            WHERE dbt_valid_from = '{snapshot_date}'
        ),
        valid_to AS (
            SELECT *,
                'before' AS record_status
            FROM {snapshot_node["relation_name"]}
            WHERE dbt_valid_to = '{snapshot_date}'
        )

        SELECT * FROM valid_to 
        UNION ALL
        SELECT * FROM valid_from;
        """

        df = dbt_adapter_query(adapter, text_query)
        return df
