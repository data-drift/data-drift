from datetime import datetime, timezone
import os
from typing import Dict, Iterator, Optional

from datagit.dataframe.helpers import (
    sort_dataframe_on_first_column_and_assert_is_unique,
)
from datagit.drift_evaluator.drift_evaluators import (
    DefaultDriftEvaluator,
    DriftEvaluatorAbstractClass,
    drift_summary_to_string,
)
from ..dataframe.dataframe_update_breakdown import (
    DataFrameUpdate,
    dataframe_update_breakdown,
)
import pandas as pd
from git import Commit, Repo


def store_table(
    *,
    store_name="default",
    table_name: str,
    table_dataframe: pd.DataFrame,
    measure_date: Optional[datetime] = None,
    drift_evaluator: DriftEvaluatorAbstractClass = DefaultDriftEvaluator(),
):
    if measure_date is None:
        measure_date = datetime.now(timezone.utc)
    table_dataframe = sort_dataframe_on_first_column_and_assert_is_unique(
        table_dataframe
    )

    local_connector = LocalConnector(store_name=store_name)

    repo = local_connector.repo
    store_dir = repo.working_dir
    print(f"Storing table {table_name} in db {store_dir}")

    latest_stored_snapshot = local_connector.get_table(table_name)

    if latest_stored_snapshot == None:
        print("Table not found, creating it")
        local_connector.init_table(
            table_name=table_name, dataframe=table_dataframe, measure_date=measure_date
        )
        print("Table stored")
        return

    else:
        update_breakdown = dataframe_update_breakdown(
            latest_stored_snapshot, table_dataframe, drift_evaluator
        )

        local_connector.handle_breakdown(
            table_name=table_name,
            update_breakdown=update_breakdown,
            measure_date=measure_date,
        )


class LocalConnector:
    home_dir = os.path.expanduser("~")
    datadrift_dir = os.path.join(home_dir, ".datadrift")
    os.makedirs(datadrift_dir, exist_ok=True)

    def __init__(self, store_name="default"):
        self.store_name = store_name
        self.repo = self.get_or_init_repo(store_name=self.store_name)
        self.store_dir = self.repo.working_dir

    def get_table(self, metric_name: str) -> Optional[pd.DataFrame]:
        store_dir = self.store_dir
        metric_file_name = f"{metric_name}.csv"
        table_file_path = os.path.join(store_dir, metric_file_name)

        if not os.path.isfile(table_file_path):
            return None
        return pd.read_csv(os.path.join(store_dir, metric_file_name))

    def init_table(
        self, table_name: str, dataframe: pd.DataFrame, measure_date: datetime
    ):
        store_dir = self.store_dir
        table_file_name = f"{table_name}.csv"
        table_file_path = os.path.join(store_dir, table_file_name)
        table_file_dir = os.path.dirname(table_file_path)
        os.makedirs(table_file_dir, exist_ok=True)
        dataframe.to_csv(table_file_path, index=True, na_rep="NA")
        add_file = [table_file_name]
        self.repo.index.add(add_file)
        self.repo.index.commit(f"NEW DATA: {table_name}", author_date=measure_date)

    def handle_breakdown(
        self,
        table_name: str,
        measure_date: datetime,
        update_breakdown: Dict[str, DataFrameUpdate],
    ):
        store_dir = self.store_dir
        table_file_name = f"{table_name}.csv"
        table_file_path = os.path.join(store_dir, table_file_name)
        for key, value in update_breakdown.items():
            value["df"].to_csv(table_file_path, na_rep="NA")
            add_file = [table_file_name]
            self.repo.index.add(add_file)
            if self.repo.index.diff("HEAD"):
                commit_message = f"{key}: {table_name}"
                if value["drift_evaluation"] != None:
                    commit_message += f"\n{value['drift_evaluation']['message']}"
                if value["drift_summary"]:
                    string_summary = drift_summary_to_string(value["drift_summary"])
                    commit_message += "\n\n" + string_summary
                self.repo.index.commit(message=commit_message, author_date=measure_date)
            else:
                pass

    def get_tables(self):
        repo = self.repo
        csv_files = [
            os.path.splitext(f)[0]
            for f in os.listdir(repo.working_dir)
            if f.endswith(".csv")
        ]
        return csv_files

    def delete_table(self, metric_name: str):
        # Getting commit history
        commit_history = list(self.get_table_history(metric_name=metric_name))

        # If there's no commit, exit
        if not commit_history:
            return

        repo = self.get_or_init_repo(store_name=self.store_name)
        active_branch = repo.active_branch

        # Create a new copy of main branch
        timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
        keep_main = f"keep_main_{timestamp}"
        repo.git.checkout("HEAD", b=keep_main)

        # Create a new temporary branch and checkout
        tmp_branch = f"tmp_branch_{timestamp}"
        repo.git.checkout("HEAD", b=tmp_branch)

        for commit in commit_history:
            print(f"Deleting commit {commit.hexsha}")
            repo.git.rebase("--onto", commit.hexsha + "^", commit.hexsha, tmp_branch)

        repo.git.branch("-D", active_branch)
        repo.git.checkout("HEAD", b=active_branch)
        repo.git.branch("-D", tmp_branch)

    def get_table_history(self, metric_name: str) -> Iterator[Commit]:
        repo = self.get_or_init_repo(store_name=self.store_name)
        metric_file_name = f"{metric_name}.csv"
        commits = repo.iter_commits(paths=metric_file_name)
        return commits

    @staticmethod
    def get_or_init_repo(store_name="default"):
        store_dir = os.path.join(LocalConnector.datadrift_dir, store_name)
        os.makedirs(store_dir, exist_ok=True)

        try:
            repo = Repo(store_dir)
            return repo
        except:
            print(f"The store {store_name} does not exist. Creating it now.")
            repo = Repo.init(store_dir)
            repo.index.commit("Init DB")
            return repo

    @staticmethod
    def delete_store(store_name="default"):
        store_dir = os.path.join(LocalConnector.datadrift_dir, store_name)
        os.system(f"rm -rf {store_dir}")
