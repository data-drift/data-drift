from io import StringIO
import os
import unittest
from datagit.dataframe.dataframe_update_breakdown import dataframe_update_breakdown
import pandas as pd


def formatDF(dict):
    df = pd.DataFrame(dict)
    df["unique_key"] = df.apply(lambda row: row["date"] + "-" + row["name"], axis=1)
    column_order = ["unique_key"] + [col for col in df.columns if col != "unique_key"]
    df = df.reindex(columns=column_order)
    return df


class TestUpdateBreakdown(unittest.TestCase):
    def setUp(self):
        base_dir = os.path.dirname(__file__)  # get the directory of the current script
        initial_csv_file = os.path.join(
            base_dir, "./datasets/test_dataframe_update_breakdown_reindexation1.csv"
        )
        self.initial_df = pd.read_csv(initial_csv_file)

        final_data = {
            "name": ["Alice", "Bob", "Charlie"],
            "date": ["2022-12", "2023-01", "2023-01"],
            "age": [25, 30, 35],
        }
        self.final_df = formatDF(final_data)

    def test_comparison_on_same_index(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        print(result)
        self.assertNotIn("Pet", result["MIGRATION Column Deleted"]["df"].columns)


if __name__ == "__main__":
    unittest.main()
