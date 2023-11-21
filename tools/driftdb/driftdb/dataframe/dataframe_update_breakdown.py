from enum import Enum
from typing import Dict, Optional, Union

import pandas as pd

from ..drift_evaluator.drift_evaluators import (
    BaseUpdateEvaluator,
    DefaultDriftEvaluator,
    DriftEvaluation,
    DriftEvaluatorContext,
    safe_drift_evaluator,
)
from ..drift_evaluator.interface import NewDataEvaluatorContext
from .helpers import reparse_dataframe
from .summarize_dataframe_updates import summarize_dataframe_updates


class UpdateType(Enum):
    DRIFT = "drift"
    OTHER = "other"


class DataFrameUpdate:
    def __init__(
        self,
        df: pd.DataFrame,
        has_update: bool,
        type: UpdateType,
        update_context: Optional[Union[DriftEvaluatorContext, NewDataEvaluatorContext]],
        update_evaluation: Optional[DriftEvaluation],
    ):
        self.df = df
        self.has_update = has_update
        self.type = type
        self.update_context = update_context
        self.update_evaluation = update_evaluation


def dataframe_update_breakdown(
    initial_dataframe: pd.DataFrame,
    final_dataframe: pd.DataFrame,
    drift_evaluator: BaseUpdateEvaluator = DefaultDriftEvaluator(),
) -> Dict[str, DataFrameUpdate]:
    if initial_dataframe.index.name != "unique_key":
        initial_dataframe = initial_dataframe.set_index("unique_key")

    if final_dataframe.index.name != "unique_key":
        final_dataframe = final_dataframe.set_index("unique_key")

    columns_added = set(final_dataframe.columns) - set(initial_dataframe.columns)
    columns_removed = set(initial_dataframe.columns) - set(final_dataframe.columns)

    step1 = initial_dataframe.drop(columns=list(columns_removed))

    # TODO handle case when there is not the collection date
    # new_data = initial_dataframe
    # if "date" in initial_dataframe.columns and "date" in final_dataframe.columns:
    new_data = final_dataframe.loc[~final_dataframe["date"].isin(initial_dataframe["date"])]

    step2 = pd.concat([step1, new_data[step1.columns]], axis=0)

    new_data_context = None
    new_data_evaluation = None
    if len(new_data) > 0:
        new_data_context = NewDataEvaluatorContext(
            before=reparse_dataframe(step1), after=reparse_dataframe(step2), added_rows=reparse_dataframe(new_data)
        )
        new_data_evaluation = drift_evaluator.compute_new_data_evaluation(new_data_context)

    step3 = final_dataframe.drop(columns=list(columns_added))

    has_drift = not step2.equals(step3)
    drift_context = None
    drift_evaluation = None
    if has_drift:
        drift_summary = summarize_dataframe_updates(initial_df=step2, final_df=step3)
        drift_context = DriftEvaluatorContext(before=step2, after=step3, summary=drift_summary)
        drift_evaluation = safe_drift_evaluator(drift_context, drift_evaluator.compute_drift_evaluation)
        # Here, in case of wrongly detected drifts, we recheck the drifts
        if (
            len(drift_summary["added_rows"]) == 0
            and len(drift_summary["deleted_rows"]) == 0
            and len(drift_summary["modified_rows_unique_keys"]) == 0
        ):
            has_drift = False

    step4 = final_dataframe.reindex(index=step3.index)

    return {
        "MIGRATION Column Deleted": DataFrameUpdate(
            df=step1,
            has_update=not initial_dataframe.equals(step1),
            type=UpdateType.OTHER,
            update_context=None,
            update_evaluation=None,
        ),
        "NEW DATA": DataFrameUpdate(
            df=step2,
            has_update=not step1.equals(step2),
            type=UpdateType.OTHER,
            update_context=new_data_context,
            update_evaluation=new_data_evaluation,
        ),
        "DRIFT": DataFrameUpdate(
            df=step3,
            has_update=has_drift,
            type=UpdateType.DRIFT,
            update_context=drift_context,
            update_evaluation=drift_evaluation,
        ),
        "MIGRATION Column Added": DataFrameUpdate(
            df=step4,
            has_update=not step3.equals(step4),
            type=UpdateType.OTHER,
            update_context=None,
            update_evaluation=None,
        ),
    }
