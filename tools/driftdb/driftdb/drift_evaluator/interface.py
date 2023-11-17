from typing import Optional, TypedDict

import pandas as pd


class DriftSummary(TypedDict):
    added_rows: pd.DataFrame
    deleted_rows: pd.DataFrame
    modified_rows_unique_keys: pd.Index
    modified_patterns: pd.DataFrame


class DriftEvaluatorContext(TypedDict):
    before: pd.DataFrame
    after: pd.DataFrame
    summary: Optional[DriftSummary]


class DriftEvaluation(TypedDict):
    should_alert: bool
    message: str
