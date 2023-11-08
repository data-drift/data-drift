from datagit.drift_evaluators import DriftEvaluatorContext
import pandas as pd
from typing import Dict, Optional, TypedDict
from enum import Enum


class UpdateType(Enum):
    DRIFT = "drift"
    OTHER = "other"


class DataFrameUpdate(TypedDict):
    df: pd.DataFrame
    has_update: bool
    type: UpdateType
    drift_context: Optional[DriftEvaluatorContext]


def dataframe_update_breakdown(
    initial_dataframe: pd.DataFrame, final_dataframe: pd.DataFrame
) -> Dict[str, DataFrameUpdate]:
    # Ensure the dataframes have the same index
    initial_dataframe = initial_dataframe.set_index(initial_dataframe.columns[0])
    final_dataframe = final_dataframe.set_index(final_dataframe.columns[0])

    # Find columns added and removed
    columns_added = set(final_dataframe.columns) - set(initial_dataframe.columns)
    columns_removed = set(initial_dataframe.columns) - set(final_dataframe.columns)

    # 1. Remove the columns
    step1 = initial_dataframe.drop(columns=columns_removed)

    # 2. Add new data based on date
    new_data = final_dataframe.loc[
        ~final_dataframe["date"].isin(initial_dataframe["date"])
    ]
    step2 = pd.concat([step1, new_data[step1.columns]], axis=0)

    step3 = final_dataframe.drop(columns=columns_added)

    # 4. Add new columns
    step4 = final_dataframe

    return {
        "MIGRATION Column Deleted": DataFrameUpdate(
            df=step1,
            has_update=not initial_dataframe.equals(step1),
            type=UpdateType.OTHER,
            drift_context=None,
        ),
        "NEW DATA": DataFrameUpdate(
            df=step2,
            has_update=not step1.equals(step2),
            type=UpdateType.OTHER,
            drift_context=None,
        ),
        "DRIFT": DataFrameUpdate(
            df=step3,
            has_update=not step2.equals(step3),
            type=UpdateType.DRIFT,
            drift_context=DriftEvaluatorContext(before=step2, after=step3),
        ),
        "MIGRATION Column Added": DataFrameUpdate(
            df=step4,
            has_update=not step3.equals(step4),
            type=UpdateType.OTHER,
            drift_context=None,
        ),
    }
