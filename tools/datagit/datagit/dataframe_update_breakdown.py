from datagit.drift_evaluators import (
    DriftEvaluation,
    DriftEvaluatorContext,
    auto_merge_drift,
    safe_drift_evaluator,
)
import pandas as pd
from typing import Callable, Dict, Optional, TypedDict
from enum import Enum


class UpdateType(Enum):
    DRIFT = "drift"
    OTHER = "other"


class DataFrameUpdate(TypedDict):
    df: pd.DataFrame
    has_update: bool
    type: UpdateType
    drift_context: Optional[DriftEvaluatorContext]
    drift_evaluation: Optional[DriftEvaluation]


def dataframe_update_breakdown(
    initial_dataframe: pd.DataFrame,
    final_dataframe: pd.DataFrame,
    drift_evaluator: Callable[
        [DriftEvaluatorContext], DriftEvaluation
    ] = auto_merge_drift,
) -> Dict[str, DataFrameUpdate]:
    if initial_dataframe.index.name != "unique_key":
        initial_dataframe = initial_dataframe.set_index("unique_key")

    if final_dataframe.index.name != "unique_key":
        final_dataframe = final_dataframe.set_index("unique_key")

    columns_added = set(final_dataframe.columns) - set(initial_dataframe.columns)
    columns_removed = set(initial_dataframe.columns) - set(final_dataframe.columns)

    step1 = initial_dataframe.drop(columns=list(columns_removed))

    # TODO handle case when there is not
    # new_data = initial_dataframe
    # if "date" in initial_dataframe.columns and "date" in final_dataframe.columns:
    new_data = final_dataframe.loc[
        ~final_dataframe["date"].isin(initial_dataframe["date"])
    ]

    step2 = pd.concat([step1, new_data[step1.columns]], axis=0)

    step3 = final_dataframe.drop(columns=list(columns_added))
    result = drift_breakdown(before_drift=step2, after_drift=step3)
    step3_1 = result["with_deleted"]
    step3_1_has_update = not step2.equals(step3_1)
    step3_1_drift_context = None
    step3_1_drift_evaluation = None
    if step3_1_has_update:
        step3_1_drift_context = DriftEvaluatorContext(before=step2, after=step3_1)
        step3_1_drift_evaluation = safe_drift_evaluator(
            step3_1_drift_context, drift_evaluator
        )
    step3_2 = result["with_deleted_and_added"]
    step3_2_has_update = not step3_1.equals(step3_2)
    step3_2_drift_context = None
    step3_2_drift_evaluation = None
    if step3_2_has_update:
        step3_2_drift_context = DriftEvaluatorContext(before=step3_1, after=step3_2)
        step3_2_drift_evaluation = safe_drift_evaluator(
            step3_2_drift_context, drift_evaluator
        )
    step3_3 = result["with_deleted_and_added_and_modified"]

    step4 = final_dataframe.reindex(index=step3_3.index, columns=step3_3.columns)

    return {
        "MIGRATION Column Deleted": DataFrameUpdate(
            df=step1,
            has_update=not initial_dataframe.equals(step1),
            type=UpdateType.OTHER,
            drift_context=None,
            drift_evaluation=None,
        ),
        "NEW DATA": DataFrameUpdate(
            df=step2,
            has_update=not step1.equals(step2),
            type=UpdateType.OTHER,
            drift_context=None,
            drift_evaluation=None,
        ),
        "DRIFT Deletion": DataFrameUpdate(
            df=step3_1,
            has_update=step3_1_has_update,
            type=UpdateType.DRIFT,
            drift_context=step3_1_drift_context,
            drift_evaluation=step3_1_drift_evaluation,
        ),
        "DRIFT Addition": DataFrameUpdate(
            df=step3_2,
            has_update=step3_2_has_update,
            type=UpdateType.DRIFT,
            drift_context=step3_2_drift_context,
            drift_evaluation=step3_2_drift_evaluation,
        ),
        "DRIFT Modification": DataFrameUpdate(
            df=step3_3,
            has_update=not step3_2.equals(step3_3),
            type=UpdateType.DRIFT,
            drift_context=DriftEvaluatorContext(before=step3_2, after=step3_3),
            drift_evaluation=None,
        ),
        "MIGRATION Column Added": DataFrameUpdate(
            df=step4,
            has_update=not step3_3.equals(step4),
            type=UpdateType.OTHER,
            drift_context=None,
            drift_evaluation=None,
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


def summarize_dataframe_updates(data_drift_context: DriftEvaluatorContext):
    """
    Summarize the updates made to a dataframe including added, deleted, and modified rows.
    Group the modifications by the pattern of changes.

    Parameters:
    - initial_df (pd.DataFrame): The original dataframe before updates.
    - final_df (pd.DataFrame): The updated dataframe after changes.
    - key (str): The name of the column or index to use as the unique key for comparison.

    Returns:
    - A dictionary with three keys: 'added', 'deleted', and 'modified', each containing
      a respective dataframe of changes, and 'modification_patterns', a dataframe summarizing
      the patterns of modification.
    """
    initial_df = data_drift_context["before"]
    final_df = data_drift_context["after"]

    initial_df = initial_df.set_index("unique_key")
    final_df = final_df.set_index("unique_key")

    deleted_rows = initial_df[~initial_df.index.isin(final_df.index)]

    added_rows = final_df[~final_df.index.isin(initial_df.index)]

    common_indices = initial_df.index.intersection(final_df.index)
    common_rows_initial = initial_df.loc[common_indices]
    common_rows_final = final_df.loc[common_indices]

    changes = common_rows_initial != common_rows_final
    changed_rows = changes[changes.any(axis=1)].index

    pattern_changes = {}
    for key in changed_rows:
        for col in common_rows_initial.columns:
            if common_rows_initial.at[key, col] != common_rows_final.at[key, col]:
                old_value = common_rows_initial.at[key, col]
                new_value = common_rows_final.at[key, col]
                change_pattern = (col, old_value, new_value)
                if change_pattern not in pattern_changes:
                    pattern_changes[change_pattern] = [key]
                else:
                    pattern_changes[change_pattern].append(key)

    patterns_list = []
    for pattern, keys in pattern_changes.items():
        col, old, new = pattern
        patterns_list.append(
            {
                "unique_keys": keys,
                "column": col,
                "old_value": old,
                "new_value": new,
                "pattern_id": hash(pattern),
            }
        )

    patterns_df = pd.DataFrame(patterns_list)

    return {
        "added": added_rows,
        "deleted": deleted_rows,
        "modified_patterns": patterns_df,
    }
