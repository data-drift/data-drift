from abc import ABC, abstractmethod
from datetime import datetime
from typing import Dict
from ..dataframe.dataframe_update_breakdown import DataFrameUpdate

import pandas as pd


class AbstractConnector(ABC):
    @abstractmethod
    def handle_breakdown(
        self,
        table_name: str,
        update_breakdown: Dict[str, DataFrameUpdate],
        measure_date: datetime,
    ):
        pass

    @abstractmethod
    def get_table(self, table_name: str) -> pd.DataFrame:
        pass

    @abstractmethod
    def init_table(
        self, table_name: str, dataframe: pd.DataFrame, measure_date: datetime
    ):
        pass
