import os

from datagit.dataset_helpers import sort_dataframe_on_first_column_and_assert_is_unique
from .dataframe_update_breakdown import dataframe_update_breakdown
import pandas as pd
from git import Repo

home_dir = os.path.expanduser("~")
datadrift_dir = os.path.join(home_dir, ".datadrift")
os.makedirs(datadrift_dir, exist_ok=True)


def store_metric(*, store_name="default", metric_name: str, metric_value: pd.DataFrame):
    metric_value = sort_dataframe_on_first_column_and_assert_is_unique(metric_value)
    store_dir = os.path.join(datadrift_dir, store_name)
    os.makedirs(store_dir, exist_ok=True)
    print(f"Storing metric {metric_name} in db {store_dir}")

    repo = Repo.init(store_dir)
    metric_file_name = f"{metric_name}.csv"
    metric_file_path = os.path.join(store_dir, metric_file_name)

    if not os.path.isfile(metric_file_path):
        metric_file_dir = os.path.dirname(metric_file_path)
        os.makedirs(metric_file_dir, exist_ok=True)
        metric_value.to_csv(metric_file_path, index=False)
        add_file = [metric_file_name]
        repo.index.add(add_file)
        repo.index.commit(f"NEW DATA: {metric_name}")
        return

    initial_dataframe = pd.read_csv(metric_file_path)
    update_breakdown = dataframe_update_breakdown(initial_dataframe, metric_value)
    for key, value in update_breakdown.items():
        value.to_csv(metric_file_path)
        add_file = [metric_file_name]
        repo.index.add(add_file)
        if repo.index.diff("HEAD"):
            repo.index.commit(f"{key}: {metric_name}")
        else:
            pass
    pass


def get_metric(*, store_name="default", metric_name: str) -> pd.DataFrame:
    store_dir = os.path.join(datadrift_dir, store_name)
    metric_file_name = f"{metric_name}.csv"
    return pd.read_csv(os.path.join(store_dir, metric_file_name))


def get_metrics(*, store_name="default"):
    store_dir = os.path.join(datadrift_dir, store_name)
    csv_files = [
        os.path.splitext(f)[0] for f in os.listdir(store_dir) if f.endswith(".csv")
    ]
    return csv_files
