import unittest
from datagit.dataframe_update_breakdown import dataframe_update_breakdown
import pandas as pd


class TestUpdateBreakdown(unittest.TestCase):
    def setUp(self):
        self.initial_df = pd.DataFrame(
            {
                "unique_key": [1, 2, 3],
                "date": ["2021-01-01", "2021-01-02", "2021-01-03"],
                "name": ["A", "B", "C"],
                "age": [25, 30, 35],
                "Pet": ["Dog", "Cat", "Bird"],
            }
        )

        self.final_df = pd.DataFrame(
            {
                "unique_key": [2, 3, 4],
                "date": ["2021-01-02", "2021-01-03", "2021-01-04"],
                "name": ["B", "C", "D"],
                "age": [30, 36, 40],
                "city": ["X", "Y", "Z"],
            }
        )

    def test_columns_deleted(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertNotIn("Pet", result["MIGRATION Column Deleted"].columns)

    def test_new_data(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertIn(4, result["NEW DATA"].index)
        self.assertEqual(result["NEW DATA"].loc[3, "age"], 35)

    def test_drift(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertNotIn(1, result["DRIFT"].index)
        self.assertEqual(result["DRIFT"].loc[3, "age"], 36)

    def test_columns_added(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertIn("city", result["MIGRATION Column Added"].columns)
        self.assertEqual(result["MIGRATION Column Added"].loc[3, "city"], "Y")


if __name__ == "__main__":
    unittest.main()
