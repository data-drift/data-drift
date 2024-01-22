from driftdb.alerting.interface import DriftSummary
from pandas import DataFrame, Index


def convert_snapshot_to_drift_summary(snapshot_diff: DataFrame, id_column="id", date_column="date") -> DriftSummary:
    required_columns = [id_column, date_column, "record_status"]
    for column in required_columns:
        if column not in snapshot_diff.columns:
            raise ValueError(f"The snapshot_diff DataFrame does not have the required column: {column}")

    initial_data = snapshot_diff[snapshot_diff["record_status"] == "before"]
    final_data = snapshot_diff[snapshot_diff["record_status"] == "after"]

    common_ids = initial_data[id_column][initial_data[id_column].isin(final_data[id_column])]

    driftSummary = DriftSummary(
        added_rows=DataFrame(),
        deleted_rows=DataFrame(),
        modified_rows_unique_keys=Index(common_ids),
        modified_patterns=DataFrame(),
    )
    return driftSummary
