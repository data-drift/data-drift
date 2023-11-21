from typing import Optional, TypedDict

import pandas as pd


class DriftSummary(TypedDict):
    added_rows: pd.DataFrame
    deleted_rows: pd.DataFrame
    modified_rows_unique_keys: pd.Index
    modified_patterns: pd.DataFrame


class DriftEvaluatorContext:
    def __init__(self, before: pd.DataFrame, after: pd.DataFrame, summary: Optional[DriftSummary]):
        self.before = before
        self.after = after
        self.summary = summary


class NewDataEvaluatorContext(TypedDict):
    before: pd.DataFrame
    after: pd.DataFrame
    added_rows: pd.DataFrame


class DriftEvaluation(TypedDict):
    should_alert: bool
    message: str
