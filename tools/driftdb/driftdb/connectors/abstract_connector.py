from abc import ABC, abstractmethod
from datetime import datetime, timezone
from typing import Dict, Optional

import pandas as pd

from ..dataframe.dataframe_update_breakdown import DataFrameUpdate, dataframe_update_breakdown
from ..dataframe.helpers import sort_dataframe_on_first_column_and_assert_is_unique
from ..drift_evaluator.drift_evaluators import BaseUpdateEvaluator, DefaultDriftEvaluator
from ..logger import get_logger
from .common import assert_valid_table_name, find_date_column, get_partition_file_path

null_logger = get_logger(__name__)


class AbstractConnector(ABC):
    logger = null_logger

    @abstractmethod
    def handle_breakdown(
        self,
        table_name: str,
        update_breakdown: Dict[str, DataFrameUpdate],
        measure_date: datetime,
    ):
        pass

    @abstractmethod
    def get_table(self, table_name: str) -> Optional[pd.DataFrame]:
        pass

    @abstractmethod
    def init_table(self, table_name: str, dataframe: pd.DataFrame, measure_date: datetime):
        pass

    def snapshot_table(
        self,
        *,
        table_dataframe: pd.DataFrame,
        table_name: str,
        measure_date: Optional[datetime] = None,
        drift_evaluator: BaseUpdateEvaluator = DefaultDriftEvaluator(),
    ):
        assert_valid_table_name(table_name)
        if measure_date is None:
            measure_date = datetime.now(timezone.utc)
        table_dataframe = sort_dataframe_on_first_column_and_assert_is_unique(table_dataframe)
        if table_dataframe.index.name != "unique_key":
            table_dataframe = table_dataframe.set_index("unique_key")
        table_dataframe = table_dataframe.astype(str)

        date_column = find_date_column(table_dataframe)
        if date_column is None:
            raise Exception("Collection date column not found")

        latest_stored_snapshot = self.get_table(table_name)

        if latest_stored_snapshot is None:
            self.logger.info("Table not found. Creating it")
            self.init_table(table_name=table_name, dataframe=table_dataframe, measure_date=measure_date)
            self.logger.info("Table stored")
            pass
        else:
            self.logger.info("Table found. Updating it")
            update_breakdown = dataframe_update_breakdown(latest_stored_snapshot, table_dataframe, drift_evaluator)
            if any(item.has_update for item in update_breakdown.values()):
                self.logger.info("Change detected")
                self.handle_breakdown(
                    table_name=table_name,
                    update_breakdown=update_breakdown,
                    measure_date=measure_date,
                )
            else:
                self.logger.info("Nothing to update")
                pass

    def partition_and_snapshot_table(
        self,
        *,
        table_dataframe: pd.DataFrame,
        measure_date: Optional[datetime] = None,
        table_name: str,
        freq: str = "M",
    ) -> None:
        """
        Partitions and snapshots a table.

        Args:
            table_dataframe (pd.DataFrame): The DataFrame representing the table to be partitioned and snapshotted.
            measure_date (Optional[datetime]): The date to be used for measuring. If not provided, the current date is used.
            table_name (str): The name of the table.
            freq (str, optional): The frequency for partitioning. Defaults to 'M' for monthly frequency. 'W' for weekly. Follows Pandas groupby date frequency alias. https://pandas.pydata.org/docs/user_guide/timeseries.html#offset-aliases

        Returns:
            None
        """
        self.logger.info("Partitionning table by month...")

        table_dataframe["date"] = pd.to_datetime(table_dataframe["date"])

        grouped = table_dataframe.groupby(pd.Grouper(key="date", freq=freq))

        # Iterate over the groups and print the sub-dataframes
        for name, group in grouped:
            self.logger.info(f"Storing table for: {name}")
            monthly_table_name = get_partition_file_path(table_name, name.strftime("%Y-%m-%d"))  # type: ignore
            self.snapshot_table(
                table_dataframe=group,
                table_name=monthly_table_name,
                measure_date=measure_date,
            )
