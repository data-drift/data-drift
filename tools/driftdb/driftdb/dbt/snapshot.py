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
    dbtRunner().invoke(["-q", "debug"], project_dir=str(project_dir))
    profile = load_profile(str(project_dir), {})
    project = load_project(str(project_dir), version_check=False, profile=profile)

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

def purge_intermediates_snapshot(snapshot_node: SnapshotNode, first_snapshot_date: str, last_snapshot_date: str):
    from dbt.adapters.factory import get_adapter
    from dbt.cli.main import dbtRunner
    from dbt.config.runtime import RuntimeConfig, load_profile, load_project

    project_dir = "."
    dbtRunner().invoke(["-q", "debug"], project_dir=str(project_dir))
    profile = load_profile(str(project_dir), {})
    project = load_project(str(project_dir), version_check=False, profile=profile)

    runtime_config = RuntimeConfig.from_parts(project, profile, {})

    adapter = get_adapter(runtime_config)

    purge_dates = []
    snapshot_dates = get_snapshot_dates(snapshot_node)
    for date in snapshot_dates[:-1]:
        if first_snapshot_date <= date <= last_snapshot_date:
            date_to_purge = date
            next_date = snapshot_dates[snapshot_dates.index(date) + 1]
            purge_dates.append({"date_to_purge": date_to_purge, "next_date": next_date})

    print("date purge", purge_dates)

    with adapter.connection_named("default"):  # type: ignore
        bookings_snapshot_purged_name = 'bookings_snapshot_purged_2'
        adapter.execute("BEGIN;")
        print("Begin transaction")
        clone_table_query = f"""
        DROP TABLE IF EXISTS {bookings_snapshot_purged_name};
        CREATE TABLE {bookings_snapshot_purged_name} AS SELECT * FROM {snapshot_node["relation_name"]};
        """
        print("Cloning table", clone_table_query)
        res = dbt_adapter_query(adapter,clone_table_query)
        
        print("Table cloned", res)

        for purge_date in purge_dates:
            date_to_purge = purge_date["date_to_purge"]
            next_date = purge_date["next_date"]
            print(f"Purging {date_to_purge} with {next_date}")
            text_query = f"""
            UPDATE      {bookings_snapshot_purged_name} SET dbt_valid_to = '{next_date}'   WHERE dbt_valid_to   = '{date_to_purge}';
            DELETE FROM {bookings_snapshot_purged_name}                                    WHERE dbt_valid_from = '{date_to_purge}' AND dbt_valid_to = '{next_date}';
            UPDATE      {bookings_snapshot_purged_name} SET dbt_valid_from = '{next_date}' WHERE dbt_valid_from = '{date_to_purge}' AND dbt_valid_to = NULL;
            """

            adapter.execute(text_query , auto_begin=True, fetch=True)
    
        adapter.execute("COMMIT;")