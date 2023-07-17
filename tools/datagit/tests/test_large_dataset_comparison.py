import os
import unittest
import pandas as pd
from datagit.github_connector import copy_and_compare_dataframes
from pandas.testing import assert_frame_equal


class TestLargeDatasetComparison(unittest.TestCase):
    def test_copy_and_compare_same_dataframe(self):
        df1 = pd.read_csv(
            os.path.join(os.path.dirname(__file__), "datasets", "ultra_large_df.csv")
        )
        df2 = pd.read_csv(
            os.path.join(os.path.dirname(__file__), "datasets", "ultra_large_df.csv")
        )

        # Call the copy_and_compare_dataframe function
        result = copy_and_compare_dataframes(df1, df2)

        # Assert that the function returns True
        if result is not None:
            self.assertEqual(len(result), 0)
        else:
            self.fail("The function did not return a dataframe")

    def test_copy_and_compare_1_line_diff_dataframe(self):
        df1 = pd.read_csv(
            os.path.join(os.path.dirname(__file__), "datasets", "ultra_large_df.csv")
        )
        df2 = pd.read_csv(
            os.path.join(os.path.dirname(__file__), "datasets", "ultra_large_df2.csv")
        )

        # Call the copy_and_compare_dataframe function
        result = copy_and_compare_dataframes(df1, df2)

        # Assert that the function returns True
        if result is not None:
            self.assertEqual(len(result), 1)
        else:
            self.fail("The function did not return a dataframe")
