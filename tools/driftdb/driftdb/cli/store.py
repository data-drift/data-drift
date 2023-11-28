import os

import pandas as pd
import typer
from driftdb.connectors.local_connector import LocalConnector

from .common import prompt_from_list

app = typer.Typer()


@app.command()
def delete(store: str = typer.Option("default", help="name of the store")):
    LocalConnector.delete_store(store_name=store)


@app.command()
def delete_table(table: str = typer.Option(None, help="name of your table")):
    local_connector = LocalConnector()
    tables = local_connector.get_tables()
    if not table:
        tables = local_connector.get_tables()
        [table, table_index] = prompt_from_list("Please enter table number", tables)
    local_connector.delete_table(table_name=table)


@app.command()
def load_csv(
    csvpathfile: str = typer.Argument(...),  # Required argument
    table: str = typer.Option(None, help="name of your table"),
    unique_key_column: str = typer.Option(None, help="name of your unique key column"),
    date_column: str = typer.Option(None, help="name of your date column"),
):
    local_connector = LocalConnector()
    if not table:
        tables = local_connector.get_tables()
        table = typer.prompt("Please enter table name (existing or not)", type=tables)

    typer.echo(f"Loading CSV file {csvpathfile}...")
    assert os.path.exists(csvpathfile), f"CSV file {csvpathfile} does not exist"
    dataframe = pd.read_csv(csvpathfile)

    if "unique_key" not in dataframe.columns:
        if not unique_key_column:
            unique_key_column = typer.prompt(
                "Please enter unique key column name", type=dataframe.columns.tolist()  # type: ignore
            )
        assert unique_key_column in dataframe.columns, f"Column {unique_key_column} does not exist in CSV file"
        dataframe.insert(0, "unique_key", dataframe[unique_key_column])

    if "date" not in dataframe.columns:
        if not date_column:
            date_column = typer.prompt(
                "Please enter date column name", type=click.Choice(dataframe.columns)  # type: ignore
            )
        assert date_column in dataframe.columns, f"Column {date_column} does not exist in CSV file"
        dataframe.insert(1, "date", dataframe[date_column])

    local_connector = LocalConnector()
    local_connector.snapshot_table(table_name=table, table_dataframe=dataframe)


if __name__ == "__main__":
    app()
