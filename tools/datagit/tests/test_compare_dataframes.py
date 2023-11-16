import unittest
from datagit.dataframe.dataframe_update_breakdown import dataframe_update_breakdown
import pandas as pd

from datagit.dataframe.helpers import generate_drift_description


class TestStoreMetric(unittest.TestCase):
    def test_compare_dataframes_with_1_addition_and_1_deletion(self):
        # Define the test dataframes
        df1 = pd.DataFrame(
            {"unique_key": [1, 2, 3], "name": ["Alice", "Bob", "Charlie"]}
        )
        df2 = pd.DataFrame(
            {"unique_key": [2, 3, 4], "name": ["Bob", "Charlie", "Dave"]}
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"]["drift_context"]
        if drift_context is None:
            raise Exception("drift_context is None")

        # Call the function being tested
        result = generate_drift_description(drift_context)

        # Define the expected result
        expected_result = "- 🆕 1 addition\n- ♻️ 0 modification\n- 🗑️ 1 deletion"

        # Assert that the actual result matches the expected result
        assert result == expected_result

    def test_compare_dataframes_with_2_deletions(self):
        # Define the test dataframes
        df1 = pd.DataFrame(
            {"unique_key": [1, 2, 3], "name": ["Alice", "Bob", "Charlie"]}
        )
        df2 = pd.DataFrame({"unique_key": [2], "name": ["Bob"]})

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"]["drift_context"]
        if drift_context is None:
            raise Exception("drift_context is None")
        # Call the function being tested
        result = generate_drift_description(drift_context)

        # Define the expected result
        expected_result = "- 🆕 0 addition\n- ♻️ 0 modification\n- 🗑️ 2 deletions"

        # Assert that the actual result matches the expected result
        assert result == expected_result

    def test_compare_dataframes_with_2_additions(self):
        # Define the test dataframes
        df1 = pd.DataFrame({"unique_key": [2], "name": ["Bob"]})
        df2 = pd.DataFrame(
            {"unique_key": [1, 2, 3], "name": ["Alice", "Bob", "Charlie"]}
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"]["drift_context"]
        if drift_context is None:
            raise Exception("drift_context is None")
        # Call the function being tested
        result = generate_drift_description(drift_context)

        # Define the expected result
        expected_result = "- 🆕 2 additions\n- ♻️ 0 modification\n- 🗑️ 0 deletion"

        # Assert that the actual result matches the expected result
        assert result == expected_result

    def test_compare_dataframes_with_2_modifications(self):
        # Define the test dataframes
        df1 = pd.DataFrame(
            {"unique_key": [1, 2, 3], "name": ["Alice", "Bob", "Charlie"]}
        )
        df2 = pd.DataFrame(
            {"unique_key": [1, 2, 3], "name": ["Alixe", "Bob", "Charles"]}
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"]["drift_context"]
        if drift_context is None:
            raise Exception("drift_context is None")
        # Call the function being tested
        result = generate_drift_description(drift_context)

        # Define the expected result
        expected_result = "- 🆕 0 addition\n- ♻️ 2 modifications\n- 🗑️ 0 deletion"

        # Assert that the actual result matches the expected result
        assert result == expected_result

    def test_compare_dataframes_with_additions_deletions_modifications(self):
        # Define the test dataframes
        df1 = pd.DataFrame(
            {"unique_key": [1, 2, 3], "name": ["Alice", "Bob", "Charlie"]}
        )
        df2 = pd.DataFrame(
            {"unique_key": [1, 3, 4], "name": ["Alixe", "Charles", "Dave"]}
        )

        break_down = dataframe_update_breakdown(df1, df2)
        drift_context = break_down["DRIFT"]["drift_context"]
        if drift_context is None:
            raise Exception("drift_context is None")
        # Call the function being tested
        result = generate_drift_description(drift_context)

        # Define the expected result
        expected_result = "- 🆕 1 addition\n- ♻️ 2 modifications\n- 🗑️ 1 deletion"

        # Assert that the actual result matches the expected result
        assert result == expected_result
