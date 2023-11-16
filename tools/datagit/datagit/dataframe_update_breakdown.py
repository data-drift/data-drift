from datagit.drift_evaluators import (
    DefaultDriftEvaluator,
    DriftEvaluation,
    DriftEvaluatorAbstractClass,
    DriftEvaluatorContext,
    DriftSummary,
    safe_drift_evaluator,
)
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
    drift_evaluation: Optional[DriftEvaluation]
    drift_summary: Optional[DriftSummary]


def dataframe_update_breakdown(
    initial_dataframe: pd.DataFrame,
    final_dataframe: pd.DataFrame,
    drift_evaluator: DriftEvaluatorAbstractClass = DefaultDriftEvaluator(),
) -> Dict[str, DataFrameUpdate]:
    if initial_dataframe.index.name != "unique_key":
        initial_dataframe = initial_dataframe.set_index("unique_key")

    if final_dataframe.index.name != "unique_key":
        final_dataframe = final_dataframe.set_index("unique_key")

    columns_added = set(final_dataframe.columns) - set(initial_dataframe.columns)
    columns_removed = set(initial_dataframe.columns) - set(final_dataframe.columns)

    step1 = initial_dataframe.drop(columns=list(columns_removed))

    # TODO handle case when there is not the collection date
    # new_data = initial_dataframe
    # if "date" in initial_dataframe.columns and "date" in final_dataframe.columns:
    new_data = final_dataframe.loc[
        ~final_dataframe["date"].isin(initial_dataframe["date"])
    ]

    step2 = pd.concat([step1, new_data[step1.columns]], axis=0)

    step3 = final_dataframe.drop(columns=list(columns_added))
    common_index = step3.index.intersection(step2.index)
    step3 = step3.reindex(index=common_index)

    has_drift = not step2.equals(step3)
    drift_summary = None
    drift_context = None
    drift_evaluation = None
    if has_drift:
        drift_summary = summarize_dataframe_updates(initial_df=step2, final_df=step3)
        drift_context = DriftEvaluatorContext(
            before=step2, after=step3, summary=drift_summary
        )
        drift_evaluation = safe_drift_evaluator(
            drift_context, drift_evaluator.compute_drift_evaluation
        )

    step4 = final_dataframe.reindex(index=step3.index)

    return {
        "MIGRATION Column Deleted": DataFrameUpdate(
            df=step1,
            has_update=not initial_dataframe.equals(step1),
            type=UpdateType.OTHER,
            drift_context=None,
            drift_evaluation=None,
            drift_summary=None,
        ),
        "NEW DATA": DataFrameUpdate(
            df=step2,
            has_update=not step1.equals(step2),
            type=UpdateType.OTHER,
            drift_context=None,
            drift_evaluation=None,
            drift_summary=None,
        ),
        "DRIFT": DataFrameUpdate(
            df=step3,
            has_update=not step2.equals(step3),
            type=UpdateType.DRIFT,
            drift_context=drift_context,
            drift_evaluation=drift_evaluation,
            drift_summary=drift_summary,
        ),
        "MIGRATION Column Added": DataFrameUpdate(
            df=step4,
            has_update=not step3.equals(step4),
            type=UpdateType.OTHER,
            drift_context=None,
            drift_evaluation=None,
            drift_summary=None,
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


def summarize_dataframe_updates(
    initial_df: pd.DataFrame,
    final_df: pd.DataFrame,
) -> DriftSummary:
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

    if initial_df.index.name != "unique_key":
        initial_df = initial_df.set_index("unique_key")

    if final_df.index.name != "unique_key":
        final_df = final_df.set_index("unique_key")

    final_df = final_df.reindex(index=initial_df.index)

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
        "added_rows": added_rows,
        "deleted_rows": deleted_rows,
        "modified_rows_unique_keys": changed_rows,
        "modified_patterns": patterns_df,
    }
