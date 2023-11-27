import typer
from driftdb.connectors.local_connector import LocalConnector

from ..dataframe.seed import generate_dataframe, insert_drift
from .common import prompt_from_list

app = typer.Typer()


@app.command()
def create(
    table: str = typer.Option(None, help="name of your table"),
    row_number: int = typer.Option(10000, help="Number of lines to generate"),
):
    if not table:
        table = typer.prompt("Table name")
    typer.echo("Creating seed file...")
    dataframe = generate_dataframe(row_number)

    typer.echo(str(dataframe.columns))
    local_connector = LocalConnector()
    local_connector.snapshot_table(table_name=table, table_dataframe=dataframe)
    typer.echo("Seed created...")


@app.command()
def update(
    table: str = typer.Option(None, help="name of your table"),
    row_number: int = typer.Option(100, help="Number of lines to update"),
):
    local_connector = LocalConnector()

    if not table:
        tables = local_connector.get_tables()
        table = prompt_from_list("Please enter table number", tables)

    typer.echo("Updating seed file...")
    dataframe = local_connector.get_table(table_name=table)
    if dataframe is None:
        raise Exception("Table not found")
    drifted_dataset = insert_drift(dataframe, row_number)

    local_connector.snapshot_table(table_name=table, table_dataframe=drifted_dataset)
    typer.echo("Seed updated...")


if __name__ == "__main__":
    app()
