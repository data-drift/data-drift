import unittest
import pandas as pd

from datagit.dataset_helpers import compare_dataframes


class TestStoreMetric(unittest.TestCase):
    def test_compare_dataframes_with_1_addition_and_1_deletion(self):
        # Define the test dataframes
        df1 = pd.DataFrame(
            {"unique_key": [1, 2, 3], "name": ["Alice", "Bob", "Charlie"]}
        )
        df2 = pd.DataFrame(
            {"unique_key": [2, 3, 4], "name": ["Bob", "Charlie", "Dave"]}
        )

        # Call the function being tested
        result = compare_dataframes(df1, df2, "unique_key")

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

        # Call the function being tested
        result = compare_dataframes(df1, df2, "unique_key")

        # Define the expected result
        expected_result = (
            "- 🆕 0 addition\n- ♻️ 0 modification\n- 🗑️ 2 deletions"
        )

        # Assert that the actual result matches the expected result
        assert result == expected_result

    def test_compare_dataframes_with_2_additions(self):
        # Define the test dataframes
        df1 = pd.DataFrame({"unique_key": [2], "name": ["Bob"]})
        df2 = pd.DataFrame(
            {"unique_key": [1, 2, 3], "name": ["Alice", "Bob", "Charlie"]}
        )

        # Call the function being tested
        result = compare_dataframes(df1, df2, "unique_key")

        # Define the expected result
        expected_result = (
            "- 🆕 2 additions\n- ♻️ 0 modification\n- 🗑️ 0 deletion"
        )

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

        # Call the function being tested
        result = compare_dataframes(df1, df2, "unique_key")

        # Define the expected result
        expected_result = (
            "- 🆕 0 addition\n- ♻️ 2 modifications\n- 🗑️ 0 deletion"
        )

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

        # Call the function being tested
        result = compare_dataframes(df1, df2, "unique_key")

        # Define the expected result
        expected_result = "- 🆕 1 addition\n- ♻️ 2 modifications\n- 🗑️ 1 deletion"

        # Assert that the actual result matches the expected result
        assert result == expected_result
