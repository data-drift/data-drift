import os
import unittest
from datetime import datetime, timedelta, timezone

import pandas as pd
from driftdb.connectors import LocalConnector


class TestLocalConnector(unittest.TestCase):
    def setUp(self):
        self.connector = LocalConnector(store_name="test")
        self.table_name = "test_table"
        self.dataframe = pd.DataFrame(
            data={
                "unique_key": ["abc", "def", "ijk"],
                "date": ["2022-01-01", "2022-01-01", "2022-01-01"],
                "some_metric": ["a", "b", "c"],
            },
        )

        self.measure_date = datetime.now(timezone.utc)

    def test_init_and_get_table(self):
        initial_dataframe = self.dataframe.copy()
        initial_dataframe = initial_dataframe.set_index("unique_key")
        self.assertIsNone(self.connector.get_table(self.table_name))

        self.connector.init_table(self.table_name, initial_dataframe, self.measure_date)

        table = self.connector.get_table(self.table_name)

        if table is None:
            self.fail("Table should be initiated")
        pd.testing.assert_frame_equal(table, self.dataframe)

    def test_snapshot_new_version(self):
        # Initialize the table
        new_dataframe = self.dataframe.copy()
        initial_dataframe = self.dataframe.copy()
        initial_dataframe = initial_dataframe.set_index("unique_key")

        self.connector.init_table(self.table_name, initial_dataframe, self.measure_date)

        new_dataframe.loc[1, "some_metric"] = "d"
        new_measure_date = self.measure_date + timedelta(days=1)
        self.connector.snapshot_table(
            table_name=self.table_name, table_dataframe=new_dataframe, measure_date=new_measure_date
        )

        # Test that the table has been updated
        table = self.connector.get_table(self.table_name)
        if table is None:
            self.fail("Table should be updated")
        pd.testing.assert_frame_equal(table, new_dataframe)

    def tearDown(self):
        # Delete the test store after each test
        self.connector.delete_store(store_name="test")


if __name__ == "__main__":
    unittest.main()
