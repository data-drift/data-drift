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
    if initial_dataframe.index.name != "unique_key":
        initial_dataframe = initial_dataframe.set_index("unique_key")

    if final_dataframe.index.name != "unique_key":
        final_dataframe = final_dataframe.set_index("unique_key")

    columns_added = set(final_dataframe.columns) - set(initial_dataframe.columns)
    columns_removed = set(initial_dataframe.columns) - set(final_dataframe.columns)

    step1 = initial_dataframe.drop(columns=columns_removed)

    # TODO handle case when there is not
    # new_data = initial_dataframe
    # if "date" in initial_dataframe.columns and "date" in final_dataframe.columns:
    new_data = final_dataframe.loc[
        ~final_dataframe["date"].isin(initial_dataframe["date"])
    ]

    step2 = pd.concat([step1, new_data[step1.columns]], axis=0)

    step3 = final_dataframe.drop(columns=columns_added)
    result = drift_breakdown(before_drift=step2, after_drift=step3)
    step3_1 = result["with_deleted"]
    step3_2 = result["with_deleted_and_added"]
    step3_3 = result["with_deleted_and_added_and_modified"]

    step4 = final_dataframe.reindex(index=step3_3.index, columns=step3_3.columns)

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
        "DRIFT Deletion": DataFrameUpdate(
            df=step3_1,
            has_update=not step2.equals(step3_1),
            type=UpdateType.DRIFT,
            drift_context=DriftEvaluatorContext(before=step2, after=step3_1),
        ),
        "DRIFT Addition": DataFrameUpdate(
            df=step3_2,
            has_update=not step3_1.equals(step3_2),
            type=UpdateType.DRIFT,
            drift_context=DriftEvaluatorContext(before=step3_1, after=step3_2),
        ),
        "DRIFT Modification": DataFrameUpdate(
            df=step3_3,
            has_update=not step3_2.equals(step3_3),
            type=UpdateType.DRIFT,
            drift_context=DriftEvaluatorContext(before=step3_2, after=step3_3),
        ),
        "MIGRATION Column Added": DataFrameUpdate(
            df=step4,
            has_update=not step3_3.equals(step4),
            type=UpdateType.OTHER,
            drift_context=None,
        ),
    }


class DriftBreakdownResult(TypedDict):
    with_deleted: pd.DataFrame
    with_deleted_and_added: pd.DataFrame
    with_deleted_and_added_and_modified: pd.DataFrame


def drift_breakdown(
    before_drift: pd.DataFrame, after_drift: pd.DataFrame
) -> DriftBreakdownResult:
    # If 'unique_key' is not the index yet, set it as the index
    if before_drift.index.name != "unique_key":
        before_drift = before_drift.set_index("unique_key")

    if after_drift.index.name != "unique_key":
        after_drift = after_drift.set_index("unique_key")

    # Find keys that were deleted, added, or stayed the same
    deleted_keys = before_drift.index.difference(after_drift.index)

    added_keys = after_drift.index.difference(before_drift.index)

    # DataFrame without deleted lines
    without_deleted = before_drift.drop(deleted_keys).sort_index()

    # DataFrame with added lines
    with_added = pd.concat([without_deleted, after_drift.loc[added_keys]]).sort_index()
    after_drift = after_drift.reindex(
        index=with_added.index, columns=with_added.columns
    )
    return DriftBreakdownResult(
        with_deleted=without_deleted,
        with_deleted_and_added=with_added,
        with_deleted_and_added_and_modified=after_drift,
    )
