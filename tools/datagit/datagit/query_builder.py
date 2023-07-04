from typing import List, Optional, Sequence


def build_query(
    table_id: str,
    unique_key_columns: List[str],
    columns: Optional[List[str]] = None,
    where_clauses: Optional[Sequence[str]] = None,
) -> str:
    if not unique_key_columns:
        raise ValueError("unique_key_columns cannot be empty")
    if where_clauses is None:
        where_clauses = ["TRUE"]
    where_clause = " AND ".join(where_clauses)
    column_str = ", ".join(columns) if columns else "*"
    unique_key_str = ", '__', ".join(unique_key_columns)

    return f"""
        SELECT
          CONCAT({unique_key_str}) AS unique_key,
          {column_str}
        FROM {table_id}
        WHERE {where_clause}
        ORDER BY 1
    """
