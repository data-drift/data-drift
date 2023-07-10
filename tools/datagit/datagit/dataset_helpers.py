import pandas as pd


def assert_dataset_has_unique_key(dataset: pd.DataFrame) -> None:
    first_column = dataset.columns[0]

    if first_column != "unique_key":
        raise ValueError(f"The first column is not named 'unique_key'")

    # Check if the first column is unique
    values = dataset[first_column].tolist()
    if len(values) != len(set(values)):
        raise ValueError(f"The {first_column} column is not unique")


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
    assert_dataset_has_date_column(df)

    df_with_parsed_dates = parse_date_column(df)

    sorted_df = df_with_parsed_dates.sort_values(by=["unique_key"])
    return sorted_df


def compare_dataframes(
    initial_df: pd.DataFrame, final_df: pd.DataFrame, unique_key: str
):
    try:
        # Get the unique keys for each dataframe
        initial_keys = set(initial_df[unique_key])
        final_keys = set(final_df[unique_key])

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

        modifications = diff.__len__()

        # Construct the result text
        result = ""
        if additions > 0:
            result += f"- 🆕 {additions} addition{'s' if additions > 1 else ''}\n"
        else:
            result += f"- ~~🆕 0 addition~~\n"
        if modifications > 0:
            result += (
                f"- ♻️ {modifications} modification{'s' if modifications > 1 else ''}\n"
            )
        else:
            result += f"- ~~♻️ 0 modification~~\n"
        if deletions > 0:
            result += f"- 🗑️ {deletions} deletion{'s' if deletions > 1 else ''}\n"
        else:
            result += f"- ~~🗑️ 0 deletion~~\n"

        return result.strip()
    except Exception as e:
        return f"Could not generate drift description: {e}"
