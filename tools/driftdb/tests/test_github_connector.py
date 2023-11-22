import unittest
from unittest import mock
from unittest.mock import MagicMock, call, patch

import pandas as pd
from driftdb.connectors.github_connector import GithubConnector


def mocked_read_csv(url, *args, **kwargs):
    return pd.DataFrame(
        {
            "unique_key": [1, 2, 3],
            "col2": [4, 5, 6],
            "date": ["2021-01-01", "2022-02-02", "2023-03-03"],
        }
    )


class TestStoreMetric(unittest.TestCase):
    def setUp(self):
        self.ghClient = MagicMock()
        self.repo = MagicMock()
        self.contents = MagicMock(download_url="url.fr")
        self.contents.decoded_content.decode.return_value = ""
        self.repo.get_contents.return_value = self.contents
        self.ghClient.get_repo.return_value = self.repo
        self.dataframe = pd.DataFrame({"unique_key": [1, 2], "col2": [3, 4], "date": ["2021-01-01", "2022-02-02"]})
        self.github_repository_name = "org/repo"
        self.table_name = "path/to/file.csv"
        self.connector = GithubConnector(
            github_client=self.ghClient,
            github_repository_name=self.github_repository_name,
            assignees=["jerome"],
        )

    def test_store_metric(self):
        with patch("pandas.read_csv", side_effect=mocked_read_csv):
            self.connector.snapshot_table(
                table_dataframe=self.dataframe,
                table_name=self.table_name,
            )

            self.repo.get_contents.assert_has_calls(
                [
                    call("path/to/file.csv", ref=self.repo.default_branch),
                    call("path/to/file.csv", ref=mock.ANY),
                ]
            )

    def test_store_metric_with_no_assignee(self):
        with patch("pandas.read_csv", side_effect=mocked_read_csv):
            self.connector.snapshot_table(
                table_dataframe=self.dataframe,
                table_name=self.table_name,
            )

            self.repo.get_contents.assert_has_calls(
                [
                    call("path/to/file.csv", ref=self.repo.default_branch),
                    call("path/to/file.csv", ref=mock.ANY),
                ]
            )

            self.repo.create_pull.assert_not_called()
