import unittest
from unittest.mock import MagicMock
from datagit.github_connector import get_valid_branch_name, create_git_branch
import re


class TestGetValidBranchName(unittest.TestCase):
    def test_valid_branch_name(self):
        filepath = "/path/to/my/file.txt"
        pattern = r"drift/\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}/path-to-my-file-txt"
        assert re.match(pattern, get_valid_branch_name(filepath))

    def test_invalid_branch_name(self):
        filepath = "/path/to/my/file with spaces.txt"
        pattern = (
            r"drift/\d{4}-\d{2}-\d{2}-\d{2}-\d{2}-\d{2}/path-to-my-file-with-spaces-txt"
        )

        assert re.match(pattern, get_valid_branch_name(filepath))


class TestGithubConnector(unittest.TestCase):
    def test_create_git_branch(self):
        mock_repo = MagicMock()
        mock_default_branch = MagicMock()
        mock_default_branch.commit.sha = "abc123"
        mock_repo.get_branch.return_value = mock_default_branch

        # Test that create_git_branch calls create_git_ref with the correct arguments
        branch_name = "my-branch"
        ref = create_git_branch(mock_repo, branch_name)
        print("ref", ref)
        mock_repo.create_git_ref.assert_any_call("refs/heads/my-branch", "abc123")


if __name__ == "__main__":
    unittest.main()
