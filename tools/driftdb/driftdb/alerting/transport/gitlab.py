from gitlab import Gitlab

from ..interface import DriftEvaluation, DriftEvaluatorContext
from .interface import AbstractAlertTransport


class GitlabAlertTransport(AbstractAlertTransport):
    def __init__(self, gitlab_client: Gitlab, project_id: str, assignees = []) -> None:
        self.gitlab_client = gitlab_client
        self.project = gitlab_client.projects.get(project_id)
        self.assignees = assignees

    def send(self, title: str, drift_evalutation: DriftEvaluation, drift_context: DriftEvaluatorContext) -> None:
        if not drift_evalutation.should_alert:
            print("No alert needed")
            return
        try:
            issue =  self.project.issues.create({'title': title, 'description': drift_evalutation.message})
            print("Issue created: " + issue.attributes['web_url'])
        except Exception as e:
            print("Error creating the issue")
            raise e