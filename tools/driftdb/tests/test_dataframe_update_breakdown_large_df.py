import os
import unittest

import pandas as pd
from driftdb.dataframe.dataframe_update_breakdown import dataframe_update_breakdown


class TestUpdateBreakdown(unittest.TestCase):
    def setUp(self):
        base_dir = os.path.dirname(__file__)  # get the directory of the current script
        initial_csv_file = os.path.join(base_dir, "./datasets/ultra_large_df.csv")
        self.initial_df = pd.read_csv(initial_csv_file)
        self.initial_df_again = pd.read_csv(initial_csv_file)

        final_csv_file = os.path.join(base_dir, "./datasets/ultra_large_df2.csv")
        self.final_df = pd.read_csv(final_csv_file)

    def test_same_df(self):
        result = dataframe_update_breakdown(self.initial_df, self.initial_df)
        all_false = all(not item.has_update for item in result.values())
        self.assertTrue(all_false)

    def test_same_df_with_different_index(self):
        result = dataframe_update_breakdown(self.initial_df, self.initial_df_again)
        all_false = all(not item.has_update for item in result.values())
        self.assertTrue(all_false)

    def test_found_drift(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertTrue(result["DRIFT"].has_update)

    def test_drift_has_1_updates(self):
        # There is only one line that has changed in large df 2
        # dfc33511-d4f6-4ddc-ba16-49d77e312282,2008-08-03,1.29,HU,Category A
        # dfc33511-d4f6-4ddc-ba16-49d77e312282,2008-08-03,1.28,HU,Category A
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        drift_context = result["DRIFT"].update_context
        if drift_context is None:
            self.fail("drift_context is None")
        summary = drift_context.summary

        self.assertIsNotNone(summary, "modified_rows_unique_keys is None")
        if summary is not None:
            self.assertEqual(len(summary["modified_rows_unique_keys"]), 1)
