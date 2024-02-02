
from ..interface import DriftEvaluation, DriftEvaluatorContext
from .interface import AbstractAlertTransport


class ConsoleAlertTransport(AbstractAlertTransport):
    def send(self, title: str, drift_evalutation: DriftEvaluation, drift_context: DriftEvaluatorContext) -> None:
        print(title)
        drift_summary = drift_context.summary
        print("added_rows \n", drift_summary["added_rows"].to_markdown())
        print("deleted_rows \n", drift_summary["deleted_rows"].to_markdown())
        print("modified_patterns \n", drift_summary["modified_patterns"].to_markdown())
        print("modified_rows_unique_keys \n", drift_summary["modified_rows_unique_keys"])

        print("should alert \n", drift_evalutation.should_alert)
        print("alert message \n", drift_evalutation.message)