import unittest
import pandas as pd
from datetime import datetime
from datagit.dataset_helpers import parse_date_column


class TestDateParsing(unittest.TestCase):

    def test_parse_date_column(self):
        # Sample input dataframe with string date formats
        data = {'date': ['2021-01-01', '2022/02/02',
                         '03-03-2023', '04.04.2024']}
        df = pd.DataFrame(data)

        # Expected output dataframe with dates in "YYYY-MM-DD" format
        expected_data = {'date': ['2021-01-01',
                                  '2022-02-02', '2023-03-03', '2024-04-04']}
        expected_df = pd.DataFrame(expected_data)

        # Apply the parse_date_column function to the input dataframe
        result_df = parse_date_column(df)

        # Assert that the output dataframe matches the expected dataframe
        self.assertTrue(result_df.equals(expected_df))

    def test_parse_date_time_column(self):
        # Sample input dataframe with string date formats
        data = {'date': ['2022-05-06 12:34:56', '2023-07-08 10:20:30',
                         '2024-09-10 09:45:00', '2025-11-12 15:30:45']}
        df = pd.DataFrame(data)

        # Expected output dataframe with dates in "YYYY-MM-DD" format
        expected_data = {'date': ['2022-05-06',
                                  '2023-07-08', '2024-09-10', '2025-11-12']}
        expected_df = pd.DataFrame(expected_data)

        # Apply the parse_date_column function to the input dataframe
        result_df = parse_date_column(df)

        # Assert that the output dataframe matches the expected dataframe
        self.assertTrue(result_df.equals(expected_df))


if __name__ == '__main__':
    unittest.main()
