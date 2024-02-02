from abc import ABC, abstractmethod

from ..interface import DriftEvaluation


class AbstractAlertTransport(ABC):
    @abstractmethod
    def send(self, title: str, drift_evalutation: DriftEvaluation) -> None:
        pass