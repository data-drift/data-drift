import traceback

import pandas as pd
from typing_extensions import Callable, List

from ..dataframe.detect_outliers import detect_outliers
from ..logger import get_logger
from .helpers import generate_drift_description
from .interface import DriftEvaluation, DriftEvaluatorContext, DriftSummary, NewDataEvaluatorContext

logger = get_logger(__name__)

DriftHandler = Callable[[DriftEvaluatorContext], DriftEvaluation]
NewDataHandler = Callable[[NewDataEvaluatorContext], DriftEvaluation]


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


def null_drift_handler(drift_context: DriftEvaluatorContext) -> DriftEvaluation:
    return DriftEvaluation(should_alert=False, message="")


def null_new_data_handler(new_data_context: NewDataEvaluatorContext) -> DriftEvaluation:
    return DriftEvaluation(should_alert=False, message="")


def DetectOutlierHandlerFactory(numerical_cols: List[str] = [], categorical_cols: List[str] = []) -> NewDataHandler:
    def handler(new_data_context: NewDataEvaluatorContext):
        outliers = detect_outliers(
            before=new_data_context.before,
            after=new_data_context.after,
            added_rows=new_data_context.added_rows,
            numerical_cols=numerical_cols,
            categorical_cols=categorical_cols,
        )
        if len(outliers) > 0:
            return DriftEvaluation(
                should_alert=True, message=f"Found {len(outliers)} outliers\n {outliers.to_markdown()}"
            )
        return DriftEvaluation(should_alert=False, message="")

    return handler


def TresholdDriftHandlerFactory(
    numerical_cols: List[str] = [],
    treshold: float = 0.1,
) -> DriftHandler:
    def handler(drift_context: DriftEvaluatorContext):
        modified_patterns = drift_context.summary["modified_patterns"]
        outliers = pd.DataFrame()
        if modified_patterns.empty:
            return DriftEvaluation(should_alert=False, message="No modified rows in drift summary")

        for col in numerical_cols:
            col_outliers = modified_patterns[
                (modified_patterns["column"] == col)
                & (modified_patterns["old_value"] != 0)
                & (
                    abs(modified_patterns["new_value"].astype(float) - modified_patterns["old_value"].astype(float))
                    / modified_patterns["old_value"].astype(float)
                    > treshold
                )
            ]
            outliers = pd.concat([outliers, col_outliers])

        if len(outliers) > 0:
            return DriftEvaluation(
                should_alert=True, message=f"Found {len(outliers)} outliers\n {outliers.to_markdown()}"
            )
        return DriftEvaluation(should_alert=False, message="")

    return handler


def alert_drift_handler(data_drift_context: DriftEvaluatorContext) -> DriftEvaluation:
    message = f"Drift detected:\n" + generate_drift_description(data_drift_context)
    return DriftEvaluation(should_alert=True, message=message)


def auto_merge_drift(data_drift_context: DriftEvaluatorContext) -> DriftEvaluation:
    message = f"Drift detected:\n" + generate_drift_description(data_drift_context)
    return DriftEvaluation(should_alert=False, message=message)


def safe_drift_evaluator(
    data_drift_context: DriftEvaluatorContext,
    drift_evaluator: Callable[[DriftEvaluatorContext], DriftEvaluation],
) -> DriftEvaluation:
    try:
        drift_evaluation = drift_evaluator(data_drift_context)
        return drift_evaluation
    except Exception as e:
        logger.warn("Drift evaluator failed: " + str(e))
        traceback.print_exc()
        logger.warn("Using default drift evaluator")
        drift_evaluation = auto_merge_drift(data_drift_context)
        return drift_evaluation
