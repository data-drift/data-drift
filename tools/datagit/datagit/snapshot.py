from sqlalchemy import Engine
from .local_connector import store_metric, get_metric_history
import pandas as pd
from datetime import datetime
import numpy as np
from tzlocal import get_localzone


def store_snapshot(
    engine: Engine, snapshot_table: str, date_column: str, unique_key: str
):
    from sqlalchemy import create_engine, text

    text_query = f"""
    SELECT DISTINCT dbt_valid_from FROM {snapshot_table}
    UNION
    SELECT DISTINCT dbt_valid_to FROM {snapshot_table} WHERE NOT NULL;
    """
    query = text(text_query)
    df = pd.read_sql(query, engine)

    metric_history = get_metric_history(metric_name=snapshot_table)
    latest_commit = next(metric_history, None)
    if latest_commit is not None:
        authored_date = latest_commit.authored_date
        authored_datetime = datetime.fromtimestamp(authored_date)
        # Compare only the second part, without the milliseconds, and take dates after the latest measure commit
        df = df.loc[df["dbt_valid_from"].dt.floor("S") > authored_datetime]

    print("Snapshot dates to process:", df)

    for index, row in df.iterrows():
        date = row["dbt_valid_from"]
        print(f"Processing data for date: {date}")
        if pd.isna(date):
            date = pd.Timestamp.now()

        asOfQuery = f"""
        SELECT * FROM {snapshot_table}
        WHERE ('{date}' >= dbt_valid_from AND  '{date}' < dbt_valid_to) OR
            (dbt_valid_to IS NULL AND dbt_valid_from <= '{date}');
        """

        # Fetch data that's valid for the given snapshot date
        data_as_of_date = pd.read_sql(text(asOfQuery), engine)

        data_as_of_date.replace({np.nan: "NA"}, inplace=True)

        # Drop dbt columns and add date and unique_key columns
        data_as_of_date = data_as_of_date.drop(
            ["dbt_scd_id", "dbt_updated_at", "dbt_valid_from", "dbt_valid_to"], axis=1
        )
        data_as_of_date["date"] = data_as_of_date[date_column]
        data_as_of_date["unique_key"] = data_as_of_date[unique_key]
        print(data_as_of_date)

        # Compute date in UTC format
        local_tz = get_localzone()
        localized_date = date.replace(tzinfo=local_tz)

        store_metric(
            metric_name=snapshot_table,
            metric_value=data_as_of_date,
            measure_date=localized_date,
        )
