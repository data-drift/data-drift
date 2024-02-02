from datetime import datetime

import inquirer
import pandas as pd
import typer
from typing_extensions import List, Optional


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

def find_date_starting_with(dates, input_date):
    for date in dates:
        if date.startswith(input_date):
            return date
    return None

def get_user_date_selection(dates: List[str], input_date: str) -> Optional[str]:
    if input_date is not None:
        if input_date == "today":
            input_date = datetime.today().date().strftime("%Y-%m-%d")
            
        
        print("dates",input_date)
        matching_date = find_date_starting_with(dates, input_date)
        return matching_date

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
