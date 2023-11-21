import unittest

import pandas as pd
from driftdb.dataframe.dataframe_update_breakdown import dataframe_update_breakdown
from driftdb.dataframe.helpers import generate_drift_description
from driftdb.drift_evaluator.interface import DriftEvaluatorContext


class TestStoreMetric(unittest.TestCase):
    def test_compare_dataframes_with_1_addition_and_1_deletion(self):
        df1 = pd.DataFrame(
            {
                "unique_key": [1, 2, 3],
                "name": ["Alice", "Bob", "Charlie"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
            }
        )
        df2 = pd.DataFrame(
            {
                "unique_key": [2, 3, 4],
                "name": ["Bob", "Charlie", "Dave"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
            }
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"].update_context
        if not isinstance(drift_context, DriftEvaluatorContext):
            raise Exception("drift_context is not the right type")

        result = generate_drift_description(drift_context)

        expected_result = "- 🆕 1 addition\n- ♻️ 0 modification\n- 🗑️ 1 deletion"

        assert result == expected_result

    def test_compare_dataframes_with_2_deletions(self):
        # Define the test dataframes
        df1 = pd.DataFrame(
            {
                "unique_key": [1, 2, 3],
                "name": ["Alice", "Bob", "Charlie"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
            }
        )
        df2 = pd.DataFrame(
            {
                "unique_key": [2],
                "name": ["Bob"],
                "date": ["2022-01-01"],
            }
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"].update_context
        if not isinstance(drift_context, DriftEvaluatorContext):
            raise Exception("drift_context is not the right type")

        result = generate_drift_description(drift_context)

        expected_result = "- 🆕 0 addition\n- ♻️ 0 modification\n- 🗑️ 2 deletions"

        assert result == expected_result

    def test_compare_dataframes_with_2_additions(self):
        df1 = pd.DataFrame(
            {
                "unique_key": [2],
                "name": ["Bob"],
                "date": ["2022-01-01"],
            }
        )
        df2 = pd.DataFrame(
            {
                "unique_key": [1, 2, 3],
                "name": ["Alice", "Bob", "Charlie"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
            }
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"].update_context
        if not isinstance(drift_context, DriftEvaluatorContext):
            raise Exception("drift_context is not the right type")

        result = generate_drift_description(drift_context)

        expected_result = "- 🆕 2 additions\n- ♻️ 0 modification\n- 🗑️ 0 deletion"

        assert result == expected_result

    def test_compare_dataframes_with_2_modifications(self):
        # Define the test dataframes
        df1 = pd.DataFrame(
            {
                "unique_key": [1, 2, 3],
                "name": ["Alice", "Bob", "Charlie"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
            }
        )
        df2 = pd.DataFrame(
            {
                "unique_key": [1, 2, 3],
                "name": ["Alixe", "Bob", "Charles"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
            }
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"].update_context
        if not isinstance(drift_context, DriftEvaluatorContext):
            raise Exception("drift_context is not the right type")

        result = generate_drift_description(drift_context)

        expected_result = "- 🆕 0 addition\n- ♻️ 2 modifications\n- 🗑️ 0 deletion"

        assert result == expected_result

    def test_compare_dataframes_with_additions_deletions_modifications(self):
        df1 = pd.DataFrame(
            {
                "unique_key": [1, 2, 3],
                "name": ["Alice", "Bob", "Charlie"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
            }
        )
        df2 = pd.DataFrame(
            {
                "unique_key": [1, 3, 4],
                "name": ["Alixe", "Charles", "Dave"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
            }
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"].update_context
        if not isinstance(drift_context, DriftEvaluatorContext):
            raise Exception("drift_context is not the right type")

        result = generate_drift_description(drift_context)

        expected_result = "- 🆕 1 addition\n- ♻️ 2 modifications\n- 🗑️ 1 deletion"

        assert result == expected_result
