from github.MainClass import Github

from .interface import AbstractAlertTransport


class GithubAlertTransport(AbstractAlertTransport):
    def __init__(self, github_client: Github) -> None:
        print("GithubAlertTransport.__init__")
        self.github_client = github_client
        super().__init__()

    def send(self, alert: str) -> None:
        raise NotImplementedError("GithubAlertTransport.send_alert is not implemented")