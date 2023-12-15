import pandas as pd
from typing import TypedDict


class DriftSummary(TypedDict):
    added_rows: pd.DataFrame
    deleted_rows: pd.DataFrame
    modified_rows_unique_keys: pd.Index
    modified_patterns: pd.DataFrame


class DriftEvaluatorContext:
    def __init__(self, before: pd.DataFrame, after: pd.DataFrame, summary: DriftSummary):
        self.before = before
        self.after = after
        self.summary = summary


class NewDataEvaluatorContext:
    def __init__(self, before: pd.DataFrame, after: pd.DataFrame, added_rows: pd.DataFrame):
        self.before = before
        self.after = after
        self.added_rows = added_rows


class DriftEvaluation:
    def __init__(self, should_alert: bool, message: str):
        self.should_alert = should_alert
        self.message = message
