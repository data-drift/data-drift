import pandas as pd


def dataframe_update_breakdown(
    initial_dataframe: pd.DataFrame, final_dataframe: pd.DataFrame
) -> dict:
    # Ensure the dataframes have the same index
    initial_dataframe = initial_dataframe.set_index(initial_dataframe.columns[0])
    final_dataframe = final_dataframe.set_index(final_dataframe.columns[0])

    # Find columns added and removed
    columns_added = set(final_dataframe.columns) - set(initial_dataframe.columns)
    columns_removed = set(initial_dataframe.columns) - set(final_dataframe.columns)

    # 1. Remove the columns
    step1 = initial_dataframe.drop(columns=columns_removed)

    # 2. Add new data based on date
    new_data = final_dataframe.loc[
        ~final_dataframe["date"].isin(initial_dataframe["date"])
    ]
    step2 = pd.concat([step1, new_data[step1.columns]], axis=0)

    step3 = final_dataframe.drop(columns=columns_added)

    # 4. Add new columns
    step4 = final_dataframe

    return {
        "MIGRATION Column Deleted": step1,
        "NEW DATA": step2,
        "DRIFT": step3,
        "MIGRATION Column Added": step4,
    }
