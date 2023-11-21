import unittest
import pandas as pd
from driftdb.dataframe.helpers import parse_date_column


class TestDateParsing(unittest.TestCase):
    def test_parse_date_time_column(self):
        data = {
            "date": [
                "2022-05-06 12:34:56",
                "2023-07-08 10:20:30",
                "2024-09-10 09:45:00",
                "2025-11-12 15:30:45",
            ]
        }
        df = pd.DataFrame(data)

        expected_data = {
            "date": ["2022-05-06", "2023-07-08", "2024-09-10", "2025-11-12"]
        }
        expected_df = pd.DataFrame(expected_data)

        result_df = parse_date_column(df)

        self.assertTrue(result_df.equals(expected_df))
