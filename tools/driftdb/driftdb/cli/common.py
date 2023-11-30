from datetime import datetime
from typing import List

import inquirer
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


def get_user_date_selection(dates: List[str]) -> str:
    page_size = 10
    page = 0

    while True:
        start = page * page_size
        end = start + page_size

        choices = dates[start:end]
        if page > 0:
            choices.insert(0, "Previous Page")  # type: ignore # navigation choice are not dates
        if end < len(dates):
            choices.append("Next Page")  # type: ignore # navigation choice are not dates

        questions = [
            inquirer.List("date", message="Choose a date", choices=choices),
        ]
        answers = inquirer.prompt(questions)

        if answers is None:
            typer.echo("No choice selected. Exiting.")
            raise typer.Exit(code=1)
        if answers["date"] == "Previous Page":
            page -= 1
            continue
        elif answers["date"] == "Next Page":
            page += 1
            continue
        else:
            return answers["date"]
