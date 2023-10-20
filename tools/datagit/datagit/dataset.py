import pandas as pd
from faker import Faker
import random
import numpy as np


def generate_dataframe(num_rows=10):
    fake = Faker()
    ids = [fake.uuid4() for _ in range(num_rows)]
    dates = [
        fake.date_between(start_date="-1y", end_date="-1m").strftime("%Y-%m-%d")
        for _ in range(num_rows)
    ]
    metric_values = [round(random.uniform(0, 10), 2) for _ in range(num_rows)]
    dataframe = pd.DataFrame(
        {"unique_key": ids, "date": dates, "metric_value": metric_values}
    )

    return dataframe


def insert_drift(dataframe: pd.DataFrame, num_drift=10):
    np.random.seed(42)
    dataframe_with_drift = dataframe.copy()

    random_indices_metric = np.random.choice(
        dataframe_with_drift.index, size=num_drift, replace=False
    )
    dataframe_with_drift.loc[random_indices_metric, "metric_value"] = [
        round(random.uniform(0, 10), 2) for _ in range(num_drift)
    ]
    return dataframe_with_drift
