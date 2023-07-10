import unittest
from datagit.query_builder import build_query


class TestBuildQuery(unittest.TestCase):
    def test_build_query(self):
        table_id = "my_table"
        columns = ["col1", "col2"]
        unique_key_columns = ["period", "organisation_id"]
        where_clauses = ["col1 > 0", "col2 < 10"]
        expected_query = """
            SELECT
              CONCAT(period, '__', organisation_id) AS unique_key,
              date AS date,
              col1, col2
            FROM my_table
            WHERE col1 > 0 AND col2 < 10
            ORDER BY 1
        """

        query = build_query(
            table_id=table_id,
            unique_key_columns=unique_key_columns,
            columns=columns,
            where_clauses=where_clauses,
            date_column="date",
        )
        self.assertEqual(
            " ".join(query.split()).strip(), " ".join(expected_query.split()).strip()
        )

    def test_build_query_no_columns(self):
        table_id = "my_table"
        columns = []
        unique_key_columns = ["period", "organisation_id"]
        where_clauses = ["col1 > 0", "col2 < 10"]
        expected_query = """
            SELECT 
              CONCAT(period, '__', organisation_id) AS unique_key,
              date AS date, 
              * 
            FROM my_table 
            WHERE col1 > 0 AND col2 < 10 
            ORDER BY 1
        """
        query = build_query(
            table_id=table_id,
            unique_key_columns=unique_key_columns,
            columns=columns,
            where_clauses=where_clauses,
            date_column="date",
        )
        self.assertEqual(
            " ".join(query.split()).strip(), " ".join(expected_query.split()).strip()
        )

    def test_build_query_no_where_clauses(self):
        table_id = "my_table"
        columns = ["col1", "col2"]
        unique_key_columns = ["period", "organisation_id"]
        expected_query = """
            SELECT 
              CONCAT(period, '__', organisation_id) AS unique_key,
              date AS date, 
              col1, col2 
            FROM my_table 
            WHERE TRUE 
            ORDER BY 1
        """
        query = build_query(
            table_id=table_id,
            unique_key_columns=unique_key_columns,
            columns=columns,
            date_column="date",
        )
        self.assertEqual(
            " ".join(query.split()).strip(), " ".join(expected_query.split()).strip()
        )

    def test_build_query_empty_unique_key_columns(self):
        table_id = "my_table"
        columns = ["col1", "col2"]
        unique_key_columns = []
        where_clauses = ["col1 > 0", "col2 < 10"]
        with self.assertRaises(ValueError):
            build_query(
                table_id=table_id,
                unique_key_columns=unique_key_columns,
                columns=columns,
                where_clauses=where_clauses,
                date_column="date",
            )
