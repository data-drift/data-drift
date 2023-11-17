from datetime import datetime
import os
from typing import Dict, Iterator, List, Optional

from .abstract_connector import AbstractConnector


from ..drift_evaluator.drift_evaluators import (
    drift_summary_to_string,
)
from ..dataframe.dataframe_update_breakdown import (
    DataFrameUpdate,
)
import pandas as pd
from git import Commit, Repo


class LocalConnector(AbstractConnector):
    home_dir = os.path.expanduser("~")
    datadrift_dir = os.path.join(home_dir, ".datadrift")
    os.makedirs(datadrift_dir, exist_ok=True)

    def __init__(self, store_name="default"):
        self.store_name = store_name
        self.repo = self.get_or_init_repo(store_name=self.store_name)
        self.store_dir = self.repo.working_dir

    def _get_table_file_path(self, table_name: str) -> str:
        table_file_name = f"{table_name}.csv"
        table_file_path = os.path.join(self.store_dir, table_file_name)
        return table_file_path

    def get_table(self, table_name: str) -> Optional[pd.DataFrame]:
        table_file_path = self._get_table_file_path(table_name)
        if not os.path.isfile(table_file_path):
            return None
        return pd.read_csv(table_file_path)

    def init_table(
        self, table_name: str, dataframe: pd.DataFrame, measure_date: datetime
    ):
        table_file_name = f"{table_name}.csv"
        table_file_path = self._get_table_file_path(table_name)

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
        table_file_name = f"{table_name}.csv"
        table_file_path = self._get_table_file_path(table_name)
        for key, value in update_breakdown.items():
            if value["has_update"]:
                print("Update: " + key)
                value["df"].to_csv(table_file_path, na_rep="NA")
                add_file = [table_file_name]
                self.repo.index.add(add_file)
                commit_message = f"{key}: {table_name}"
                if value["drift_evaluation"] != None:
                    commit_message += f"\n{value['drift_evaluation']['message']}"
                if value["drift_summary"]:
                    string_summary = drift_summary_to_string(value["drift_summary"])
                    commit_message += "\n\n" + string_summary
                self.repo.index.commit(message=commit_message, author_date=measure_date)

    def get_tables(self):
        repo = self.repo
        csv_files = [
            os.path.splitext(f)[0]
            for f in os.listdir(repo.working_dir)
            if f.endswith(".csv")
        ]
        return csv_files

    def delete_table(self, table_name: str):
        # Getting commit history
        commit_history = list(self.get_table_history(table_name=table_name))

        # If there's no commit, exit
        if not commit_history:
            return

        repo = self.repo
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

    def get_table_history(self, table_name: str) -> Iterator[Commit]:
        table_name = f"{table_name}.csv"
        commits = self.repo.iter_commits(paths=table_name)
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
