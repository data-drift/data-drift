import datetime
import os
import re


def get_alert_branch_name(filepath: str) -> str:
    """
    Returns a valid Git branch name based on the given filepath.
    """
    # Replace any non-alphanumeric characters with hyphens
    branch_name = re.sub(r"[^a-zA-Z0-9]+", "-", filepath)

    # Remove any leading or trailing hyphens
    branch_name = branch_name.strip("-")

    # Convert to lowercase
    branch_name = branch_name.lower()
    now = datetime.datetime.now()
    datetime_str = now.strftime("%Y-%m-%d-%H-%M-%S")
    # Append a prefix
    branch_name = f"drift/{datetime_str}/{branch_name}"

    # Truncate to 63 characters (the maximum allowed length for a Git branch name)
    branch_name = branch_name[:63]

    return branch_name


def find_date_column(df):
    date_columns = df.filter(like="date").columns
    if len(date_columns) > 0:
        return date_columns[0]
    else:
        return df.columns[0]


def get_partition_file_path(file_path, month):
    directory, file_name = os.path.split(file_path)
    file_name, extension = os.path.splitext(file_name)

    new_file_name = f"{file_name}/{month}{extension}"

    new_file_path = os.path.join(directory, new_file_name)

    return new_file_path


def assert_valid_table_name(table_name: str):
    if table_name.startswith("/"):
        raise ValueError("Table name cannot start with a '/'")
