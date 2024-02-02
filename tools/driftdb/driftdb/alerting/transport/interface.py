from abc import ABC, abstractmethod

from ..interface import DriftEvaluation, DriftEvaluatorContext


class AbstractAlertTransport(ABC):
    @abstractmethod
    def send(self, title: str, drift_evalutation: DriftEvaluation, drift_context: DriftEvaluatorContext) -> None:
        pass