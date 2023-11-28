import typer
from driftdb.dbt.snapshot import get_snapshot_nodes

app = typer.Typer()


@app.command()
def show():
    snapshot_nodes = get_snapshot_nodes()
    print(snapshot_nodes)
