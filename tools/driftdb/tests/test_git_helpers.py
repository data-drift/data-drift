import unittest
from unittest.mock import MagicMock
from driftdb.connectors.github_connector import get_alert_branch_name
import re


class TestGetValidBranchName(unittest.TestCase):
    def test_valid_branch_name(self):
        filepath = "/path/to/my/file.txt"
        pattern = r"drift/\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}/path-to-my-file-txt"
        assert re.match(pattern, get_alert_branch_name(filepath))

    def test_invalid_branch_name(self):
        filepath = "/path/to/my/file with spaces.txt"
        pattern = (
            r"drift/\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}/path-to-my-file-with-spaces-txt"
        )

        assert re.match(pattern, get_alert_branch_name(filepath))


if __name__ == "__main__":
    unittest.main()
