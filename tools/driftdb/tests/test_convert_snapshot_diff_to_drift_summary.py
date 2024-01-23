import unittest

import pandas as pd
import pytest
from driftdb.dbt.snapshot_to_drift import convert_snapshot_to_drift_summary


class TestConvertSnapshotDiffToDriftSummary(unittest.TestCase):
    def test_convert_empty_snapshot(self):
        with pytest.raises(ValueError):
            convert_snapshot_to_drift_summary(pd.DataFrame())

    def test_convert_snapshot_with_1_modification(self):
        data = {
            "month": [pd.Timestamp("2023-02-01"), pd.Timestamp("2023-02-01")],
            "monthly_metric": [22547.9, 22547.9],
            "dbt_scd_id": ["534dc75b6f6e75f74e2f6358959a3582", "db8cd259cd772f56da85969038a4267f"],
            "dbt_updated_at": [pd.Timestamp("2023-12-26 09:23:18.569704"), pd.Timestamp("2024-01-22 16:01:31.618521")],
            "dbt_valid_from": [pd.Timestamp("2023-12-26 09:23:18.569704"), pd.Timestamp("2024-01-22 16:01:31.618521")],
            "dbt_valid_to": [pd.Timestamp("2024-01-22 16:01:31.618521"), pd.NaT],
            "record_status": ["before", "after"],
        }

        df = pd.DataFrame(data)
        context = convert_snapshot_to_drift_summary(snapshot_diff=df, id_column="month", date_column="month")
        assert context.summary["added_rows"].empty
        assert context.summary["deleted_rows"].empty
        assert context.summary["modified_rows_unique_keys"].equals(pd.Index([pd.Timestamp("2023-02-01")]))

    def test_convert_snapshot_with_1_addition(self):
        data = {
            "month": [pd.Timestamp("2022-11-01")],
            "monthly_metric": [4007.95],
            "dbt_scd_id": ["98353a44eca48d39bbb466c5351255b9"],
            "dbt_updated_at": [pd.Timestamp("2023-12-26 09:23:18.569704")],
            "dbt_valid_from": [pd.Timestamp("2023-12-26 09:23:18.569704")],
            "dbt_valid_to": [pd.NaT],
            "record_status": ["after"],
        }

        df = pd.DataFrame(data)

        context = convert_snapshot_to_drift_summary(snapshot_diff=df, id_column="month", date_column="month")

        assert context.summary["added_rows"].equals(
            pd.DataFrame(
                {
                    "month": [pd.Timestamp("2022-11-01")],
                    "monthly_metric": [4007.95],
                }
            )
        )

    def test_convert_snapshot_with_1_deletion(self):
        data = {
            "month": [pd.Timestamp("2022-11-01")],
            "monthly_metric": [4007.95],
            "dbt_scd_id": ["98353a44eca48d39bbb466c5351255b9"],
            "dbt_updated_at": [pd.Timestamp("2023-12-26 09:23:18.569704")],
            "dbt_valid_from": [pd.Timestamp("2023-12-26 09:23:18.569704")],
            "dbt_valid_to": [pd.NaT],
            "record_status": ["before"],
        }

        df = pd.DataFrame(data)

        context = convert_snapshot_to_drift_summary(snapshot_diff=df, id_column="month", date_column="month")

        assert context.summary["deleted_rows"].equals(
            pd.DataFrame(
                {
                    "month": [pd.Timestamp("2022-11-01")],
                    "monthly_metric": [4007.95],
                }
            )
        )
