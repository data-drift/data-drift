import unittest
import pandas as pd
from datagit.dataset_helpers import parse_date_column


class TestDateParsing(unittest.TestCase):

    # def test_parse_date_column(self):
    #     data1 = {'date': ['2021-01-01', '2022/02/02',
    #                       '03-03-2023', '04.04.2024']}
    #     df1 = pd.DataFrame(data1)

    #     expected_data1 = {'date': ['2021-01-01',
    #                                '2022-02-02', '2023-03-03', '2024-04-04']}
    #     expected_df1 = pd.DataFrame(expected_data1)

    #     result_df1 = parse_date_column(df1)

    #     self.assertTrue(result_df1.equals(expected_df1))

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
