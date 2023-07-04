import base64
import unittest
from unittest.mock import MagicMock, call
import pandas as pd
from github import GithubException
from datagit.github_connector import store_metric

csv_content = '''
unique_key,date
1,3
2,4
'''


class TestStoreMetric(unittest.TestCase):
    def setUp(self):
        self.ghClient = MagicMock()
        self.repo = MagicMock()
        self.contents = MagicMock(
            content=base64.b64encode(bytes(csv_content, 'utf-8')))
        self.contents.decoded_content.decode.return_value = ""
        self.repo.get_contents.return_value = self.contents
        self.ghClient.get_repo.return_value = self.repo
        self.dataframe = pd.DataFrame({"unique_key": [1, 2], "col2": [3, 4]})
        self.filepath = "org/repo/path/to/file.csv"

    def test_store_metric(self):
        store_metric(self.ghClient, self.dataframe,
                     self.filepath, ["jerome"], "production", False)
        self.repo.get_contents.assert_has_calls(
            [
                call("path/to/file.csv", ref=self.repo.default_branch),
                call("path/to/file.csv", ref="production"),
            ]
        )

    def test_store_metric_pull_request_already_exists(self):
        self.repo.create_pull.side_effect = GithubException(
            422, {"message": "A pull request already exists"}, None)
        store_metric(self.ghClient, self.dataframe,
                     self.filepath, ["jerome"], "production", False)

        self.repo.get_contents.assert_has_calls(
            [
                call("path/to/file.csv", ref=self.repo.default_branch),
                call("path/to/file.csv", ref="production"),
            ]
        )

    def test_store_metric_with_no_assignee(self):
        store_metric(self.ghClient, self.dataframe,
                     self.filepath, [], "production", False)

        self.repo.get_contents.assert_has_calls(
            [
                call("path/to/file.csv", ref=self.repo.default_branch),
                call("path/to/file.csv", ref="production"),
            ]
        )

        self.repo.create_pull.assert_not_called()
