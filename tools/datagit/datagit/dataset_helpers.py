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
    df.loc[duplicates, "unique_key"] = (
        df["unique_key"][duplicates] + "-duplicate-" + counter.astype(str)
    )

    return df


def convert_object_to_string(df):
    # Check each column
    for col in df.columns:
        # If column data type is 'object', convert it to 'string'
        if df[col].dtype == "object":
            df[col] = df[col].astype("string")
    return df


def compare_dataframes(
    initial_df: pd.DataFrame, final_df: pd.DataFrame, unique_key: str
):
    try:
        additions, deletions, diff = get_addition_deletion_and_diff(
            initial_df, final_df, unique_key
        )

        modifications = diff.__len__()

        # Construct the result text
        result = ""
        if additions > 0:
            result += f"- 🆕 {additions} addition{'s' if additions > 1 else ''}\n"
        else:
            result += f"- 🆕 0 addition\n"
        if modifications > 0:
            result += (
                f"- ♻️ {modifications} modification{'s' if modifications > 1 else ''}\n"
            )
        else:
            result += f"- ♻️ 0 modification\n"
        if deletions > 0:
            result += f"- 🗑️ {deletions} deletion{'s' if deletions > 1 else ''}\n"
        else:
            result += f"- 🗑️ 0 deletion\n"

        return result.strip()

    except Exception as e:
        return f"Could not generate drift description: {e}"


def get_addition_deletion_and_diff(initial_df, final_df, unique_key):
    if not unique_key in initial_df.columns:
        initial_df = initial_df.reset_index()
    if not unique_key in final_df.columns:
        final_df = final_df.reset_index()

    initial_df = convert_object_to_string(initial_df)
    final_df = convert_object_to_string(final_df)

    # Get the unique keys for each dataframe
    initial_keys = set(initial_df.reset_index()[unique_key])
    final_keys = set(final_df.reset_index()[unique_key])

    intersection_keys = initial_keys.intersection(final_keys)

    # Calculate the additions, modifications, and deletions
    additions = len(final_keys - intersection_keys)
    deletions = len(initial_keys - intersection_keys)

    # Get the intersection of the unique keys

    # Filter the rows in df1 and df2 that match the intersection of the unique keys
    df1_intersection = initial_df[
        initial_df[unique_key].isin(intersection_keys)
    ].reset_index(drop=True)
    df2_intersection = final_df[
        final_df[unique_key].isin(intersection_keys)
    ].reset_index(drop=True)

    diff = df1_intersection.compare(df2_intersection)
    return additions, deletions, diff
