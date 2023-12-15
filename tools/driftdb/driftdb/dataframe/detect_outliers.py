from typing import Optional

import pandas as pd


def detect_outliers(
    before: pd.DataFrame,
    after: pd.DataFrame,
    added_rows: pd.DataFrame,
    numerical_cols: list[str],
    categorical_cols: list[str],
):
    old_df = before
    new_lines = added_rows
    outliers = pd.DataFrame()

    for col in numerical_cols:
        old_col = pd.to_numeric(old_df[col], errors="coerce")
        Q1 = old_col.quantile(0.25)
        Q3 = old_col.quantile(0.75)
        IQR = Q3 - Q1

        lower_bound = Q1 - 1.5 * IQR
        upper_bound = Q3 + 1.5 * IQR

        new_col = pd.to_numeric(new_lines[col], errors="coerce")
        is_outlier = (new_col < lower_bound) | (new_col > upper_bound)
        col_outliers = new_lines[is_outlier].copy()
        col_outliers["Reason"] = f"Column {col} out of boundaries"
        outliers = pd.concat([outliers, col_outliers])

    for col in categorical_cols:
        if col == "unique_key":
            continue
        if col == "date":
            continue
        old_categories = set(old_df[col].unique())

        new_categories = set(new_lines[col].unique()) - old_categories
        is_new_category = new_lines[col].isin(new_categories)
        cat_outliers = new_lines[is_new_category].copy()
        cat_outliers["Reason"] = f"Column {col} new unkown category"

        outliers = pd.concat([outliers, cat_outliers])

    # Drop duplicate rows
    outliers = outliers.drop_duplicates()

    return outliers
