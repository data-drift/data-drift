import pandas as pd
from driftdb.drift_evaluator.drift_evaluators import DriftSummary


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

    initial_df = initial_df.astype(str)
    final_df = final_df.astype(str)

    deleted_rows = initial_df[~initial_df.index.isin(final_df.index)]

    added_rows = final_df[~final_df.index.isin(initial_df.index)]

    common_indices = initial_df.index.intersection(final_df.index)
    common_rows_initial = initial_df.loc[common_indices]
    common_rows_final = final_df.loc[common_indices]
    common_rows_final = common_rows_final.reindex(index=common_rows_initial.index)

    changes = common_rows_initial != common_rows_final
    changed_rows_index = changes[changes.any(axis=1)].index

    # There may be rows that have not changed but pandas will consider them as changed
    pattern_changes = {}
    for key in changed_rows_index:
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
        "modified_rows_unique_keys": changed_rows_index,
        "modified_patterns": patterns_df,
    }
