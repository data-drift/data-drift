from datetime import datetime, timezone
import os
from typing import Iterator

from datagit.dataset_helpers import sort_dataframe_on_first_column_and_assert_is_unique
from .dataframe_update_breakdown import dataframe_update_breakdown
import pandas as pd
from git import Commit, Repo

home_dir = os.path.expanduser("~")
datadrift_dir = os.path.join(home_dir, ".datadrift")
os.makedirs(datadrift_dir, exist_ok=True)


def store_metric(
    *,
    store_name="default",
    metric_name: str,
    metric_value: pd.DataFrame,
    measure_date=datetime.now(timezone.utc),
):
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
        metric_value.to_csv(metric_file_path, index=False, na_rep="NA")
        add_file = [metric_file_name]
        repo.index.add(add_file)
        repo.index.commit(f"NEW DATA: {metric_name}", author_date=measure_date)
        return

    initial_dataframe = pd.read_csv(metric_file_path)
    update_breakdown = dataframe_update_breakdown(initial_dataframe, metric_value)
    for key, value in update_breakdown.items():
        value.to_csv(metric_file_path, na_rep="NA")
        add_file = [metric_file_name]
        repo.index.add(add_file)
        if repo.index.diff("HEAD"):
            repo.index.commit(f"{key}: {metric_name}", author_date=measure_date)
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


def delete_metric(*, store_name="default", metric_name: str):
    # Getting commit history
    commit_history = list(
        get_metric_history(store_name=store_name, metric_name=metric_name)
    )

    # If there's no commit, exit
    if not commit_history:
        return

    store_dir = os.path.join(datadrift_dir, store_name)
    repo = Repo(store_dir)

    # Create a new copy of main branch
    timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
    keep_main = f"keep_main_{timestamp}"
    repo.git.checkout("HEAD", b=keep_main)

    # Create a new temporary branch and checkout
    repo.git.checkout("HEAD", b="tmp_branch")

    for commit in commit_history:
        repo.git.rebase("--onto", commit.hexsha + "^", commit.hexsha, "tmp_branch")

    repo.git.branch("-D", "main")
    repo.git.checkout("HEAD", b="main")
    repo.git.branch("-D", "tmp_branch")


def get_metric_history(*, store_name="default", metric_name: str) -> Iterator[Commit]:
    store_dir = os.path.join(datadrift_dir, store_name)
    metric_file_name = f"{metric_name}.csv"
    repo = Repo(store_dir)
    commits = repo.iter_commits(paths=metric_file_name)
    return commits
