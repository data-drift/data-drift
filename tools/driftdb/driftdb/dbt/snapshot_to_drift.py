import traceback

from pandas import DataFrame, Index

from ..alerting.interface import DriftEvaluatorContext, DriftSummary
from ..logger import get_logger

logger = get_logger("summarize_dataframe_updates")


def convert_snapshot_to_drift_summary(
    snapshot_diff: DataFrame, id_column="id", date_column="date"
) -> DriftEvaluatorContext:
    required_columns = [id_column, date_column, "record_status"]
    for column in required_columns:
        if column not in snapshot_diff.columns:
            logger.warn(
                f"The snapshot_diff DataFrame does not have the required column: {column}. Select a different id_column or date_column from {snapshot_diff.columns}"
            )
            raise ValueError(f"The snapshot_diff DataFrame does not have the required column: {column}")

    initial_data = snapshot_diff[snapshot_diff["record_status"] == "before"].drop(
        columns=["dbt_scd_id", "dbt_updated_at", "dbt_valid_from", "dbt_valid_to", "record_status"]
    )
    final_data = snapshot_diff[snapshot_diff["record_status"] == "after"].drop(
        columns=["dbt_scd_id", "dbt_updated_at", "dbt_valid_from", "dbt_valid_to", "record_status"]
    )

    common_ids = initial_data[id_column][initial_data[id_column].isin(final_data[id_column])]
    added_rows = final_data[~final_data[id_column].isin(common_ids)]
    deleted_rows = initial_data[~initial_data[id_column].isin(common_ids)]

    initial_data.set_index(id_column, inplace=True)
    final_data.set_index(id_column, inplace=True)

    # There may be rows that have not changed but pandas will consider them as changed
    pattern_changes = {}
    for key in common_ids:
        for col in initial_data.columns:
            try:
                if initial_data.at[key, col] != final_data.at[key, col]:
                    old_value = initial_data.at[key, col]
                    new_value = final_data.at[key, col]
                    change_pattern = (col, old_value, new_value)
                    if change_pattern not in pattern_changes:
                        pattern_changes[change_pattern] = [key]
                    else:
                        pattern_changes[change_pattern].append(key)
            except:
                logger.warn(
                    f"Error while processing pattern change in row {key} and column {col} \n {traceback.format_exc()}"
                )

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

    patterns_df = DataFrame(patterns_list)

    driftSummary = DriftSummary(
        added_rows=added_rows,
        deleted_rows=deleted_rows,
        modified_rows_unique_keys=Index(common_ids),
        modified_patterns=DataFrame(patterns_df),
    )
    return DriftEvaluatorContext(initial_data, final_data, driftSummary)
