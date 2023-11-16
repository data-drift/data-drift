from datetime import datetime, timezone
from typing import Optional
from .local_connector import LocalConnector
from .common import find_date_column
from .github_connector import GithubConnector
from ..dataframe.dataframe_update_breakdown import (
    dataframe_update_breakdown,
)
from ..dataframe.helpers import (
    sort_dataframe_on_first_column_and_assert_is_unique,
)
from ..drift_evaluator.drift_evaluators import (
    DefaultDriftEvaluator,
    DriftEvaluatorAbstractClass,
)

import pandas as pd


def snapshot_table_local(
    table_dataframe: pd.DataFrame,
    table_name: str,
    connector: LocalConnector,
    drift_evaluator: DriftEvaluatorAbstractClass = DefaultDriftEvaluator(),
    measure_date: Optional[datetime] = None,
):
    if measure_date is None:
        measure_date = datetime.now(timezone.utc)
    table_dataframe = sort_dataframe_on_first_column_and_assert_is_unique(
        table_dataframe
    )
    if table_dataframe.index.name != "unique_key":
        table_dataframe = table_dataframe.set_index("unique_key")
    table_dataframe = table_dataframe.astype("string")
    date_column = find_date_column(table_dataframe)
    if date_column is None:
        raise Exception("Collection date column not found")

    latest_stored_snapshot = connector.get_table(table_name)

    if latest_stored_snapshot is None:
        print("Table not found, creating it")
        connector.init_table(
            table_name=table_name, dataframe=table_dataframe, measure_date=measure_date
        )
        print("Table stored")
        pass
    else:
        print("Table found, updating it")
        update_breakdown = dataframe_update_breakdown(
            latest_stored_snapshot, table_dataframe, drift_evaluator
        )
        if any(item["has_update"] for item in update_breakdown.values()):
            print("Change detected")
            connector.handle_breakdown(
                table_name=table_name,
                update_breakdown=update_breakdown,
                measure_date=measure_date,
            )
        else:
            print("Nothing to update")
            pass


def snapshot_table(
    table_dataframe: pd.DataFrame,
    table_name: str,
    github_connector: GithubConnector,
    drift_evaluator: DriftEvaluatorAbstractClass = DefaultDriftEvaluator(),
):
    table_dataframe = sort_dataframe_on_first_column_and_assert_is_unique(
        table_dataframe
    )
    if table_dataframe.index.name != "unique_key":
        table_dataframe = table_dataframe.set_index("unique_key")
    table_dataframe = table_dataframe.astype("string")
    date_column = find_date_column(table_dataframe)
    if date_column is None:
        raise Exception("Collection date column not found")

    latest_stored_snapshot = github_connector.get_latest_table_snapshot(table_name)

    if latest_stored_snapshot is None:
        print("Table not found, creating it")
        github_connector.init_file(file_path=table_name, dataframe=table_dataframe)
        print("Table stored")
        pass
    else:
        print("Table found, updating it")
        update_breakdown = dataframe_update_breakdown(
            latest_stored_snapshot, table_dataframe, drift_evaluator
        )
        if any(item["has_update"] for item in update_breakdown.values()):
            print("Change detected")
            github_connector.handle_breakdown(
                table_name=table_name, update_breakdown=update_breakdown
            )
        else:
            print("Nothing to update")
            pass


def partition_and_snapshot_table(
    *,
    github_connector: GithubConnector,
    table_dataframe: pd.DataFrame,
    table_name: str,
) -> None:
    print("Partitionning table by month...")

    table_dataframe["date"] = pd.to_datetime(table_dataframe["date"])

    grouped = table_dataframe.groupby(pd.Grouper(key="date", freq="M"))

    # Iterate over the groups and print the sub-dataframes
    for name, group in grouped:
        print(f"Storing table for Month: {name}")
        monthly_table_name = get_monthly_file_path(table_name, name.strftime("%Y-%m"))  # type: ignore
        snapshot_table(
            table_dataframe=group,
            table_name=monthly_table_name,
            github_connector=github_connector,
        )
