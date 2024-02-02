from github import GithubException, Repository
from github.MainClass import Github

from ..interface import DriftEvaluation, DriftEvaluatorContext
from .interface import AbstractAlertTransport


class GithubAlertTransport(AbstractAlertTransport):
    def __init__(self, github_client: Github, repository_name: str, assignees = []) -> None:
        self.github_client = github_client
        self.repo = github_client.get_repo(repository_name)
        self.assignees = assignees

    def send(self, title: str, drift_evalutation: DriftEvaluation, drift_context: DriftEvaluatorContext) -> None:
        if not drift_evalutation.should_alert:
            print("No alert needed")
            return
        try:
            issue = self.repo.create_issue(
                title=title,
                body=drift_evalutation.message,
                assignees=self.assignees,
            )
            print("Issue created: " + issue.html_url)
        except GithubException as e:
            if e.status == 422:
                print("Issue already exists. skipping...")
            else:
                raise e