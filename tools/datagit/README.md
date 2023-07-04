# Datagit

**Datagit** is a git based metric store

```python
>>> from datagit import github_connector

>>> dataframe = bigquery.Client().query(query).to_dataframe()
{"unique_key": ['2022-01-01_FR', '2022-01-01_GB'...
>>> github_connector.store_metric(Github("Token"), dataframe=dataframe, filepath="Samox/datagit/data/act_metrics_finance/mrr.csv", assignees=["Samox"])
'🎉 data/act_metrics_finance/mrr.csv Successfully stored!'
'💩 Historical data change detected, Samox was assigned to it'
```

# Getting Started

To get started with Datagit, follow these steps:

1. Create a new repository on GitHub called `datagit` with a README file.
2. Generate a personal access token on GitHub that has access to the `datagit` repository. You can do this by going to your GitHub settings, selecting "Developer settings", and then "Personal access tokens". Click "Generate new token" and give it the necessary permissions (content and pull requests).
3. In your data pipelines, when relevant, call `store_metric` with the following parameters
   - a github client with your token Github("Token")
   - your metric in a dataframe format
   - the path of metric in a with a csv format: "your_orga/your_repo/path/to/your.csv"
   - The owner of the metric

For instance

```python
>>> from datagit import github_connector
>>> github_connector.store_metric(Github("Token"), dataframe=dataframe, filename="Samox/datagit/data/act_metrics_finance/mrr.csv", assignee=["Samox"])
```

That's it! With these steps, you can start using Datagit to store and track your metrics over time.

## Example

```python
>>> githubToken = "github_pat****"
>>> githubRepo = "ReplaceOrgaName/ReplaceRepoName"
>>> import pandas as pd
>>> dataframe = pd.DataFrame({'unique_key': ['a', 'b', 'c'], 'amount': [1001, 1002, 1003], 'is_active': [True, False, True]})
>>> from github import Github
>>> from datagit import github_connector
>>> github_connector.store_metric(Github(githubToken), dataframe=dataframe, filename=githubRepo+"data/act_metrics_finance/mrr.csv")
```

# Dataframe

Datagit is base on the standard dataframe format from [Pandas](https://pandas.pydata.org/docs/).
One can use any library to get the data as long as the format fits the following requirements:

1. The first column of the dataframe must be `unique_key`
2. The first columns must have only unique keys

## Query Builder

Datagit provides a simple query builder to store a table:

```python
>>> from datagit import query_builder
>>> query = query_builder.build_query(table_id="my_table", unique_key_columns=["organisation_id", "date_month"])
'SELECT CONCAT(organisation_id, '__', date_month) AS unique_key, * FROM my_table WHERE TRUE ORDER BY 1'
```

More [examples here](tests/test_query_builder.py)
