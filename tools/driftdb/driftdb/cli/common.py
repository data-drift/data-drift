import pandas as pd
import typer


def prompt_from_list(prompt: str, choices: list):
    for idx, choice in enumerate(choices, start=1):
        typer.echo(f"{idx}: {choice}")

    while True:
        user_input = typer.prompt(prompt)
        if user_input.isdigit() and 1 <= int(user_input) <= len(choices):
            return [choices[int(user_input) - 1], int(user_input) - 1]
        else:
            typer.echo("Invalid input. Please enter a number from the list.")


def dbt_adapter_query(
    adapter,
    query: str,
) -> pd.DataFrame:
    _, table = adapter.execute(query, fetch=True)
    data = {column_name: table.columns[column_name].values() for column_name in table.column_names}
    return pd.DataFrame(data)
