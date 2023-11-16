import traceback
from typing import Callable, Optional, TypedDict

from .interface import DriftEvaluation, DriftEvaluatorContext, DriftSummary
from ..dataframe.helpers import generate_drift_description
import pandas as pd
from abc import ABC, abstractmethod


def drift_summary_to_string(drift_summary: DriftSummary) -> str:
    return (
        f"Drift Summary:\n"
        f"Added Rows:\n{drift_summary['added_rows'].to_string()}\n"
        f"Deleted Rows:\n{drift_summary['deleted_rows'].to_string()}\n"
        f"Modified Rows Unique Keys:\n{drift_summary['modified_rows_unique_keys']}\n"
        f"Modified Patterns:\n{drift_summary['modified_patterns'].to_string()}\n"
        f"Drift Summary Json Begin\n"
        f"Added Rows JSON:\n{drift_summary['added_rows'].to_json(index=True)}\n"
        f"Deleted Rows JSON:\n{drift_summary['deleted_rows'].to_json(index=True)}\n"
        f"Modified Rows Unique Keys List:\n{drift_summary['modified_rows_unique_keys'].to_list()}\n"
        f"Modified Patterns Json:\n{drift_summary['modified_patterns'].to_json()}\n"
        f"Drift Summary Json End\n"
    )


def parse_drift_summary(commit_message: str) -> DriftSummary:
    # Extracting the drift summary part
    start_tag = "Drift Summary Json Begin\n"
    end_tag = "\nDrift Summary Json End"

    if start_tag not in commit_message or end_tag not in commit_message:
        raise ValueError("Drift summary tags not found in the commit message.")

    start_index = commit_message.find(start_tag) + len(start_tag)
    end_index = commit_message.find(end_tag)
    drift_summary_str = commit_message[start_index:end_index]

    # Parsing each part back into DataFrames and Index
    parts = drift_summary_str.split("\n")
    drift_summary = DriftSummary(
        added_rows=pd.read_json(parts[1]),
        deleted_rows=pd.read_json(parts[3]),
        modified_rows_unique_keys=pd.Index(parts[5].split(",")),
        modified_patterns=pd.read_json(parts[7]),
    )

    return drift_summary


class DriftEvaluatorAbstractClass(ABC):
    @staticmethod
    @abstractmethod
    def compute_drift_evaluation(
        data_drift_context: DriftEvaluatorContext,
    ) -> DriftEvaluation:
        pass


class DefaultDriftEvaluator(DriftEvaluatorAbstractClass):
    @staticmethod
    def compute_drift_evaluation(
        data_drift_context: DriftEvaluatorContext,
    ) -> DriftEvaluation:
        return auto_merge_drift(data_drift_context)


class AlertDriftEvaluator(DriftEvaluatorAbstractClass):
    @staticmethod
    def compute_drift_evaluation(
        data_drift_context: DriftEvaluatorContext,
    ) -> DriftEvaluation:
        return alert_drift(data_drift_context)


def alert_drift(data_drift_context: DriftEvaluatorContext) -> DriftEvaluation:
    message = f"Drift detected:\n" + generate_drift_description(data_drift_context)
    return {"should_alert": True, "message": message}


def auto_merge_drift(data_drift_context: DriftEvaluatorContext) -> DriftEvaluation:
    message = f"Drift detected:\n" + generate_drift_description(data_drift_context)
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
