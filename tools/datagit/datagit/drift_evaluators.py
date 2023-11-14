import traceback
from typing import Callable, TypedDict
from datagit.dataset_helpers import compare_dataframes
import pandas as pd


class DriftSummary(TypedDict):
    added_rows: pd.DataFrame
    deleted_rows: pd.DataFrame
    modified_rows_unique_keys: pd.Index
    modified_patterns: pd.DataFrame


def drift_summary_to_string(drift_summary: DriftSummary) -> str:
    return (
        f"Drift Summary:\n"
        f"Added Rows:\n{drift_summary['added_rows'].to_string()}\n"
        f"Deleted Rows:\n{drift_summary['deleted_rows'].to_string()}\n"
        f"Modified Rows Unique Keys:\n{drift_summary['modified_rows_unique_keys']}\n"
        f"Modified Patterns:\n{drift_summary['modified_patterns'].to_string()}\n"
        f"Drift Summary Json\n"
        f"Added Rows JSON:\n{drift_summary['added_rows'].to_json(index=True)}\n"
        f"Deleted Rows JSON:\n{drift_summary['deleted_rows'].to_json(index=True)}\n"
        f"Modified Rows Unique Keys List:\n{drift_summary['modified_rows_unique_keys'].to_list()}\n"
        f"Modified Patterns Json:\n{drift_summary['modified_patterns'].to_json()}\n"
        f"Drift Summary Json End\n"
    )


class DriftEvaluatorContext(TypedDict):
    before: pd.DataFrame
    after: pd.DataFrame
    summary: DriftSummary


class DriftEvaluation(TypedDict):
    should_alert: bool
    message: str


DriftEvaluator = Callable[[DriftEvaluatorContext], DriftEvaluation]


def alert_drift(data_drift_context: DriftEvaluatorContext) -> DriftEvaluation:
    message = f"Drift detected:\n" + compare_dataframes(
        data_drift_context["before"],
        data_drift_context["after"],
        "unique_key",
    )
    return {"should_alert": True, "message": message}


def auto_merge_drift(data_drift_context: DriftEvaluatorContext) -> DriftEvaluation:
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
    drift_evaluator: Callable[[DriftEvaluatorContext], DriftEvaluation],
) -> DriftEvaluation:
    try:
        drift_evaluation = drift_evaluator(data_drift_context)
        return drift_evaluation
    except Exception as e:
        print("Drift evaluator failed: " + str(e))
        traceback.print_exc()
        print("Using default drift evaluator")
        drift_evaluation = auto_merge_drift(data_drift_context)
        return drift_evaluation
