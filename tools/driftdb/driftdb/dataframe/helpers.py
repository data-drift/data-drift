import pandas as pd


def assert_dataset_has_unique_key(dataset: pd.DataFrame) -> None:
    first_column = dataset.columns[0]

    if first_column != "unique_key":
        raise ValueError(f"The first column is not named 'unique_key'")


def assert_dataset_column_is_unique(dataset: pd.DataFrame, column_name: str) -> None:
    values = dataset[column_name].tolist()
    if len(values) != len(set(values)):
        raise ValueError(f"The {column_name} column is not unique")


def assert_dataset_has_date_column(dataset: pd.DataFrame) -> None:
    if "date" not in dataset.columns:
        raise ValueError(f"The dataset does not have a date column")


def parse_date_column(dataset: pd.DataFrame) -> pd.DataFrame:
    column = "date"
    date_values = pd.to_datetime(dataset[column], errors="coerce")

    formatted_dates = date_values.dt.strftime("%Y-%m-%d")

    dataset[column] = formatted_dates
    return dataset


def sort_dataframe_on_first_column_and_assert_is_unique(
    df: pd.DataFrame,
) -> pd.DataFrame:
    assert_dataset_has_unique_key(df)
    df["unique_key"] = df["unique_key"].astype(str)
    df = rename_duplicates(df)
    assert_dataset_column_is_unique(df, "unique_key")
    assert_dataset_has_date_column(df)

    df_with_parsed_dates = parse_date_column(df)

    sorted_df = df_with_parsed_dates.sort_values(by=["unique_key"])
    return sorted_df


def rename_duplicates(df):
    """
    Rename duplicate 'unique_key' values in the DataFrame.

    Parameters:
    - df: pandas DataFrame with a 'unique_key' column.

    Returns:
    - DataFrame with renamed duplicates.
    """

    # Find duplicated rows based on the 'unique_key' column
    duplicates = df["unique_key"].duplicated(keep="first")

    # Create a series with the same index as the dataframe for counting duplicates
    counter = df[duplicates].groupby("unique_key").cumcount() + 1

    # Format the 'unique_key' for duplicates
    df.loc[duplicates, "unique_key"] = df["unique_key"][duplicates] + "-duplicate-" + counter.astype(str)

    return df


def convert_object_to_string(df):
    # Check each column
    for col in df.columns:
        # If column data type is 'object', convert it to 'string'
        if df[col].dtype == "object":
            df[col] = df[col].astype("string")
    return df


def reparse_dataframe(df):
    df = df.copy()
    for col in df.columns:
        # Try to convert to numeric
        df[col] = pd.to_numeric(df[col], errors="ignore")

        # If the column is still object type, try to convert to datetime
        if df[col].dtype == "object":
            df[col] = pd.to_datetime(df[col], errors="ignore")

        # Check for potential categorical data
        if df[col].dtype == "object":
            num_unique_values = len(df[col].unique())
            num_total_values = len(df[col])
            if num_unique_values / num_total_values < 0.2:  # Threshold for categorical data
                df[col] = df[col].astype("category")

    return df
