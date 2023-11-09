import traceback
from typing import Callable, TypedDict
from datagit.dataset_helpers import compare_dataframes
import pandas as pd


class DriftEvaluatorContext(TypedDict):
    before: pd.DataFrame
    after: pd.DataFrame


def alert_drift(data_drift_context: DriftEvaluatorContext):
    message = f"Drift detected:\n" + compare_dataframes(
        data_drift_context["before"],
        data_drift_context["after"],
        "unique_key",
    )
    return {"should_alert": True, "message": message}


def auto_merge_drift(data_drift_context: DriftEvaluatorContext):
    message = f"Drift detected:\n" + compare_dataframes(
        data_drift_context["before"],
        data_drift_context["after"],
        "unique_key",
    )
    return {
        "should_alert": False,
        "message": message,
    }


def safe_drift_evaluator(
    data_drift_context: DriftEvaluatorContext,
    drift_evaluator: Callable[[DriftEvaluatorContext], dict],
):
    try:
        drift_evaluation = drift_evaluator(data_drift_context)
        return drift_evaluation
    except Exception as e:
        print("Drift evaluator failed: " + str(e))
        traceback.print_exc()
        print("Using default drift evaluator")
        drift_evaluation = auto_merge_drift(data_drift_context)
        return drift_evaluation


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
    # Implementation goes here...
    # Similar to the code block provided in the previous message

    initial_df = data_drift_context["before"]
    final_df = data_drift_context["after"]
    # Assuming initial_df and final_df are pre-loaded pandas DataFrames

    # If 'unique_key' is not the index yet, set it as the index
    initial_df = initial_df.set_index("unique_key")
    final_df = final_df.set_index("unique_key")

    # Step 1: Identify deleted rows
    deleted_rows = initial_df[~initial_df.index.isin(final_df.index)]

    # Step 2: Identify added rows
    added_rows = final_df[~final_df.index.isin(initial_df.index)]

    # Step 3: Identify modified rows and group by patterns of modification
    # Find rows with the same 'unique_key' in both dataframes
    common_indices = initial_df.index.intersection(final_df.index)
    common_rows_initial = initial_df.loc[common_indices]
    common_rows_final = final_df.loc[common_indices]

    # Detect changes
    changes = common_rows_initial != common_rows_final
    changed_rows = changes[changes.any(axis=1)].index

    # Initialize a dictionary to hold patterns of changes
    pattern_changes = {}

    # Iterate over changed rows and record the pattern of change
    for key in changed_rows:
        change_pattern = []
        for col in common_rows_initial.columns:
            if common_rows_initial.at[key, col] != common_rows_final.at[key, col]:
                old_value = common_rows_initial.at[key, col]
                new_value = common_rows_final.at[key, col]
                change_pattern.append((col, old_value, new_value))

        # Convert the list of changes to a hashable tuple
        change_pattern = tuple(change_pattern)

        # Group by change pattern
        if change_pattern not in pattern_changes:
            pattern_changes[change_pattern] = [key]
        else:
            pattern_changes[change_pattern].append(key)

    # Convert the pattern_changes to a more structured form, such as a DataFrame
    patterns_list = []
    for pattern, keys in pattern_changes.items():
        for change in pattern:
            col, old, new = change
            patterns_list.append(
                {
                    "unique_keys": keys,
                    "column": col,
                    "old_value": old,
                    "new_value": new,
                    "pattern_id": hash(pattern),  # Unique identifier for the pattern
                }
            )

    # Convert list of patterns to DataFrame
    patterns_df = pd.DataFrame(patterns_list)

    return {
        "added": added_rows,
        "deleted": deleted_rows,
        "modified_patterns": patterns_df,
    }
