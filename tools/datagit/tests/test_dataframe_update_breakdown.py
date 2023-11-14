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
                "country": ["FR", "SP", "NA"],
                "age": [25, 30, 35],
                "Pet": ["Dog", "Cat", "Bird"],
            }
        )

        self.final_df = pd.DataFrame(
            {
                "unique_key": [2, 3, 4],
                "date": ["2021-01-02", "2021-01-03", "2021-01-04"],
                "name": ["B", "C", "D"],
                "country": ["SP", "NA", "BE"],
                "age": [30, 36, 40],
                "city": ["X", "Y", "Z"],
            }
        )

    def test_columns_deleted(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertNotIn("Pet", result["MIGRATION Column Deleted"]["df"].columns)

    def test_new_data(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertIn(4, result["NEW DATA"]["df"].index)
        self.assertEqual(result["NEW DATA"]["df"].loc[3, "age"], 35)

    def test_drift(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertNotIn(1, result["DRIFT"]["df"].index)
        self.assertEqual(result["DRIFT"]["df"].loc[3, "age"], 36)

    def test_columns_added(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertIn("city", result["MIGRATION Column Added"]["df"].columns)
        self.assertEqual(result["MIGRATION Column Added"]["df"].loc[3, "city"], "Y")

    def test_country_code_NA(self):
        result = dataframe_update_breakdown(self.initial_df, self.final_df)
        self.assertFalse(
            result["MIGRATION Column Deleted"]["df"]["country"].isna().any()
        )
        self.assertFalse(result["NEW DATA"]["df"]["country"].isna().any())
        self.assertFalse(result["DRIFT"]["df"]["country"].isna().any())
        self.assertFalse(result["MIGRATION Column Added"]["df"]["country"].isna().any())
        self.assertEqual(
            result["MIGRATION Column Deleted"]["df"].loc[3, "country"], "NA"
        )
        self.assertEqual(result["MIGRATION Column Added"]["df"].loc[3, "country"], "NA")


if __name__ == "__main__":
    unittest.main()
