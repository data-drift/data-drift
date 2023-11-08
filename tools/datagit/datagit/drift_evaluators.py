from typing import TypedDict
from datagit.dataset_helpers import compare_dataframes
import pandas as pd


class DriftEvaluatorContext(TypedDict):
    before: pd.DataFrame
    after: pd.DataFrame


def default_drift_evaluator(data_drift_context: DriftEvaluatorContext):
    alert_message = f"Drift detected:\n" + compare_dataframes(
        data_drift_context["before"],
        data_drift_context["after"],
        "unique_key",
    )
    return {"should_alert": True, "message": alert_message}


def auto_merge_drift(data_drift_context: DriftEvaluatorContext):
    return {
        "should_alert": False,
        "message": "Drift detected and automatically merged.",
    }
