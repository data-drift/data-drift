import traceback
from typing import Callable, TypedDict
from datagit.dataset_helpers import compare_dataframes
import pandas as pd


class DriftEvaluatorContext(TypedDict):
    before: pd.DataFrame
    after: pd.DataFrame


def alert_drift(data_drift_context: DriftEvaluatorContext):
    message = f"Drift detected:\n" + compare_dataframes(
        data_drift_context["before"],
        data_drift_context["after"],
        "unique_key",
    )
    return {"should_alert": True, "message": message}


def auto_merge_drift(data_drift_context: DriftEvaluatorContext):
    message = f"Drift detected:\n" + compare_dataframes(
        data_drift_context["before"],
        data_drift_context["after"],
        "unique_key",
    )
    return {
        "should_alert": False,
        "message": message,
    }


def safe_drift_evaluator(
    data_drift_context: DriftEvaluatorContext,
    drift_evaluator: Callable[[DriftEvaluatorContext], dict],
):
    try:
        drift_evaluation = drift_evaluator(data_drift_context)
        return drift_evaluation
    except Exception as e:
        print("Drift evaluator failed: " + str(e))
        traceback.print_exc()
        print("Using default drift evaluator")
        drift_evaluation = auto_merge_drift(data_drift_context)
        return drift_evaluation
