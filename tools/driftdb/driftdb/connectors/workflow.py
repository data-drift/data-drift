from datetime import datetime, timezone
from typing import Optional

from driftdb.connectors.abstract_connector import AbstractConnector
from .common import find_date_column
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


def snapshot_table(
    *,
    connector: AbstractConnector,
    table_dataframe: pd.DataFrame,
    table_name: str,
    measure_date: Optional[datetime] = None,
    drift_evaluator: DriftEvaluatorAbstractClass = DefaultDriftEvaluator(),
):
    if measure_date is None:
        measure_date = datetime.now(timezone.utc)
    table_dataframe = sort_dataframe_on_first_column_and_assert_is_unique(
        table_dataframe
    )
    if table_dataframe.index.name != "unique_key":
        table_dataframe = table_dataframe.set_index("unique_key")
    table_dataframe = table_dataframe.astype(str)

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


def partition_and_snapshot_table(
    *,
    connector: AbstractConnector,
    table_dataframe: pd.DataFrame,
    measure_date: Optional[datetime] = None,
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
            measure_date=measure_date,
            connector=connector,
        )
