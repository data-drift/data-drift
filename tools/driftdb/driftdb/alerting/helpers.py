import traceback

from .interface import DriftEvaluatorContext


def generate_drift_description(drift_context: DriftEvaluatorContext):
    if drift_context.summary is None:
        return f"Could not generate drift description"
    try:
        summary = drift_context.summary
        additions = len(summary["added_rows"])
        deletions = len(summary["deleted_rows"])

        modifications = len(summary["modified_rows_unique_keys"])

        # Construct the result text
        result = ""
        if additions > 0:
            result += f"- 🆕 {additions} addition{'s' if additions > 1 else ''}\n"
        else:
            result += f"- 🆕 0 addition\n"
        if modifications > 0:
            result += f"- ♻️ {modifications} modification{'s' if modifications > 1 else ''}\n"
        else:
            result += f"- ♻️ 0 modification\n"
        if deletions > 0:
            result += f"- 🗑️ {deletions} deletion{'s' if deletions > 1 else ''}\n"
        else:
            result += f"- 🗑️ 0 deletion\n"

        return result.strip()

    except Exception as e:
        traceback.print_exc()
        return f"Could not generate drift description: {e}"
