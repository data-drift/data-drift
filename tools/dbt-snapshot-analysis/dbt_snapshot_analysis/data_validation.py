import pandas as pd


def assert_valid_columns(df: pd.DataFrame):
    assert 'dbt_valid_from' in df.columns, "Column 'dbt_valid_from' not found in DataFrame"
    assert 'dbt_valid_to' in df.columns, "Column 'dbt_valid_to' not found in DataFrame"


def assert_column_is_date(df: pd.DataFrame, col_name: str):
    """
    Asserts that a column in a pandas DataFrame is convertible to a datetime format.

    Parameters:
    df (pandas.DataFrame): The DataFrame to check.
    col_name (str): The name of the column to check.

    Raises:
    AssertionError: If the column is not convertible to a datetime format.
    """
    try:
        pd.to_datetime(df[col_name], errors='raise', unit='s')
    except ValueError as e:
        error_msg = f"Column '{col_name}' is not convertible to a datetime format. Error message: {str(e)}"
        raise AssertionError(error_msg) from e
