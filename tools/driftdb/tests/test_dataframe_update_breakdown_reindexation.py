import os
import unittest

import pandas as pd
from driftdb.dataframe.dataframe_update_breakdown import dataframe_update_breakdown
from driftdb.dataframe.helpers import sort_dataframe_on_first_column_and_assert_is_unique
from driftdb.drift_evaluator.interface import DriftEvaluatorContext


def formatDF(dict):
    df = pd.DataFrame(dict)
    df["unique_key"] = df.apply(lambda row: row["date"] + "-" + row["name"], axis=1)
    column_order = ["unique_key"] + [col for col in df.columns if col != "unique_key"]
    df = df.reindex(columns=column_order)
    return df


class TestUpdateBreakdown(unittest.TestCase):
    def setUp(self):
        base_dir = os.path.dirname(__file__)  # get the directory of the current script
        initial_csv_file = os.path.join(base_dir, "./datasets/test_dataframe_update_breakdown_reindexation1.csv")
        self.initial_df = pd.read_csv(
            initial_csv_file,
            dtype="string",
            keep_default_na=False,
        )

        final_data = {
            "name": ["Alice", "Bob", "Charlie"],
            "date": ["2022-12", "2023-01", "2023-01"],
            "age": [25, 30, 35],
        }
        self.final_df = sort_dataframe_on_first_column_and_assert_is_unique(formatDF(final_data))

    def test_comparison_on_same_index(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        drift_context = result["DRIFT"].update_context
        if not isinstance(drift_context, DriftEvaluatorContext):
            raise Exception("drift_context is not the right type")

        modified_rows_unique_keys = drift_context.summary

        self.assertIsNotNone(modified_rows_unique_keys, "modified_rows_unique_keys is None")
        if modified_rows_unique_keys is not None:
            self.assertEqual(len(modified_rows_unique_keys["modified_rows_unique_keys"]), 1)
            self.assertEqual(
                modified_rows_unique_keys["modified_rows_unique_keys"][0],
                "2023-01-Charlie",
            )


if __name__ == "__main__":
    unittest.main()
