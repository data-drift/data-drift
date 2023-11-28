import json
from datetime import datetime
from typing import List, TypedDict

import pandas as pd

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


def get_snapshot_dates(snapshot_node: SnapshotNode) -> List[datetime]:
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
        ORDER BY dbt_valid_from;
        """

        df = dbt_adapter_query(adapter, text_query)
        df["dbt_valid_from"] = pd.to_datetime(df["dbt_valid_from"])
        return df["dbt_valid_from"].tolist()
