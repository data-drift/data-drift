import time
from typing import Optional, List
import pandas as pd
from github import Github, Repository, ContentFile, GithubException


class GithubConnector:
    def __init__(
        self,
        github_client: Github,
        github_repository_name: str,
        default_branch: Optional[str] = None,
        assignees: Optional[List[str]] = None,
    ):
        self.repo = github_client.get_repo(github_repository_name)
        self.default_branch = (
            default_branch if default_branch is not None else self.repo.default_branch
        )
        self.assert_branch_exist(self.repo, self.default_branch)
        self.assignees = assignees if assignees is not None else []

    def assert_file_exists(self, file_path: str) -> Optional[ContentFile.ContentFile]:
        try:
            contents = self.repo.get_contents(file_path, ref=self.default_branch)
            assert not isinstance(contents, list), "pathfile returned multiple contents"
            return contents
        except GithubException as e:
            if e.status == 404:
                return None
            else:
                raise e

    def get_latest_table_snapshot(self, table_name: str) -> Optional[pd.DataFrame]:
        table_file_content = self.assert_file_exists(table_name)
        if table_file_content is None:
            return None
        else:
            old_dataframe = pd.read_csv(
                table_file_content.download_url,
                dtype="string",
                keep_default_na=False,
            )
            return old_dataframe

    def init_file(
        self,
        file_path: str,
        dataframe: pd.DataFrame,
    ):
        commit_message = "New data: " + file_path
        print("Commit: " + commit_message)
        self.repo.create_file(
            file_path,
            commit_message,
            dataframe.to_csv(index=True, header=True),
            self.default_branch,
        )

    def create_pullrequest(self, file_path: str, description_body: str, branch: str):
        try:
            if len(self.assignees) > 0:
                pullrequest = self.repo.create_pull(
                    title="New drift detected " + file_path,
                    body=description_body,
                    head=branch,
                    base=self.default_branch,
                )
                print("Pull request created: " + pullrequest.html_url)
                existing_assignees = self.assert_assignees_exists()
                pullrequest.add_to_assignees(*existing_assignees)
            else:
                print("No assignees, skipping pull request creation")
        except GithubException as e:
            if e.status == 422:
                print("Pull request already exists, skipping...")
            else:
                raise e

    @staticmethod
    def assert_branch_exist(repo: Repository.Repository, branch_name: str) -> None:
        try:
            branch = repo.get_branch(branch_name)
        except:
            branch = None

        if not branch:
            print(f"Branch {branch_name} doesn't exist, creating it...")
            reported_branch = repo.get_branch(repo.default_branch)

            repo.create_git_ref(f"refs/heads/{branch_name}", reported_branch.commit.sha)

    def assert_assignees_exists(self) -> List[str]:
        members = [collaborator.login for collaborator in self.repo.get_collaborators()]
        exising_assignees = []
        for assignee in self.assignees:
            if assignee not in members:
                print(f"Assignee {assignee} does not exist")
            else:
                exising_assignees.append(assignee)
        return exising_assignees

    def checkout_branch_from_default_branch(self, branch_name: str):
        self.assert_branch_exist(self.repo, branch_name)
        try:
            ref = self.repo.get_git_ref(f"heads/{branch_name}")
            ref.delete()
        except GithubException:
            pass

        """
        Checkout a branch from the default branch of the given repository.
        """
        # Get the default branch of the repository
        default_branch = self.repo.get_branch(self.default_branch)
        print(
            "Checkout branch: " + branch_name,
            " from branch:" + default_branch.name,
        )

        # Create a new reference to the default branch

        try:
            ref = self.repo.create_git_ref(
                f"refs/heads/{branch_name}", default_branch.commit.sha
            )
        except GithubException:
            pass
        return

    def update_file_with_retry(
        self,
        file_path,
        commit_message,
        data,
        branch,
        max_retries=3,
    ):
        retries = 0

        while retries < max_retries:
            try:
                content = self.assert_file_exists(file_path)
                if content is None:
                    response = self.repo.create_file(
                        file_path, commit_message, data, branch
                    )
                    print(response["commit"].html_url)
                else:
                    response = self.repo.update_file(
                        file_path, commit_message, data, content.sha, branch
                    )
                    print(response["commit"].html_url)
                return
            except GithubException as e:
                if e.status == 409:
                    retries += 1
                    time.sleep(1)  # Wait for 1 second before retrying
                else:
                    raise e
        raise Exception(f"Failed to update file after {max_retries} retries")
