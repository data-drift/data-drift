from datagit.dataset_helpers import compare_dataframes


def default_drift_evaluator(data_drift_context):
    alert_message = f"Drift detected:\n" + compare_dataframes(
        data_drift_context["reported_dataframe"],
        data_drift_context["computed_dataframe"],
        "unique_key",
    )
    return {"should_alert": True, "message": alert_message}


def auto_merge_drift(data_drift_context):
    return {
        "should_alert": False,
        "message": "Drift detected and automatically merged.",
    }
