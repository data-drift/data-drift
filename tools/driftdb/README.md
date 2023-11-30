# Driftdb

<p align="center">
  <a href="https://pypi.org/project/drift/">
    <img src="https://img.shields.io/pypi/v/drift?style=flat-square" alt="DriftDb version">
  </a>
  <a href="https://pypi.org/project/driftdb/">
    <img src="https://img.shields.io/pypi/dm/driftdb?style=flat-square" alt="DriftDb monthly downloads">
  </a>
</p>

**Driftdb** is a historical metric store

```python
from driftdb.connectors import GithubConnector
from github import Github

github_connector = GithubConnector(github_client=Github("gh_token"), github_repository_name="org/repo")

dataframe = bigquery.Client().query(query).to_dataframe()
{"unique_key": ['2022-01-01_FR', '2022-01-01_GB'...

github_connector.snapshot_table(table_dataframe=dataframe, table_name="revenue")
'🎉 data/act_metrics_finance/mrr.csv Successfully stored!'
'💩 Historical data change detected, Ammy was assigned to it'
```

# Purpose

Non-moving data is a journey, in reality, the data moves and it has many impacts (Data Integrity and Reconciliation, Predictive Modeling, Historical Data Accuracy)
The purpose of this library is:

- to snapshot the data, and parse the diff in chunks (schema update, new data collection, data duplication, drift...)
- to store it using a connector
- to trigger alerts

# Getting Started (with Github as a store, it's free)

To get started with Driftdb, follow these steps:

1. Create a new repository on GitHub called `datadrift` (or whatever other name you prefer) with a README file.
2. Generate a personal access token on GitHub that has access to the `datadrift` repository. You can do this by going to your GitHub settings, selecting "Developer settings", and then "Personal access tokens". Click "Generate new token" and give it the necessary permissions (content and pull requests).
3. In your data pipelines, when relevant, call `snapshot_table` with the following parameters
   - a connector (in this example a github connector)
   - your table in a dataframe format
   - the name of the table: "kpi/my_kpi"

For instance

```python
>>> from driftdb.connectors import GithubConnector
>>> github_connector = GithubConnector(github_client=Github("gh_token"), github_repository_name="org/repo")
>>> github_connector.snapshot_table(table_dataframe=dataframe, table_name="revenue")
```

That's it! With these steps, you can start using Driftdb to store and track your metrics over time.

# Dataframe

Driftdb is base on the standard dataframe format from [Pandas](https://pandas.pydata.org/docs/).
One can use any library to get the data as long as the format fits the following requirements:

1. The first column of the dataframe must be `unique_key`
2. The first columns must have only unique keys
3. The second column must be a date (which is the collection date: the booking_date, the order_date etc)

The granularity of the dataframe depends on every use case:

- it can be at very low level (like transaction) or aggregated (like a metric)
- it can contain all the dimension, or none

## 1st column: Unique key

The unique_key is used to detect a modification in historical data

In case you have duplicated lines, driftdb will automatically rename them with `-duplicate-n`

```plaintext
  unique_key  value
0          A     10
1          B     20
2          C     30
3          B     40
4          C     50
5          C     60
6          D     70
```

```
         unique_key  value
0                A     10
1                B     20
2                C     30
3    B-duplicate-1     40
4    C-duplicate-1     50
5    C-duplicate-2     60
6                D     70
```

## 2nd column: Date

The date key is used to detect new historical data, or deleted historical data. And differentiate if a new batch is being collected (which won't be a drift)

# Large Dataset

## Partitionning

In case of more than 1M rows, partitionning is recomanded using the `partition_and_store_table` function.

```python
>>> from driftdb.connectors.workflow import partition_and_snapshot_table

>>> very_large_dataframe = bigquery.Client().query(query).to_dataframe()
{"unique_key": ['2022-01-01_FR', '2022-01-01_GB'...
>>> connector.partition_and_snapshot_table(table_dataframe=very_large_dataframe, table_name="act_metrics_finance/mrr")
'🎁 Partitionning data/act_metrics_finance/mrr.csv...'
```

# Drift

A drift is a modification of historical data. It can be a modification, addition or deletion in a table that is supposed to be "non-moving data".

## Drift Evaluator

A drift evaluator is a class that implement the following abstract class:

```python
class BaseDriftEvaluator(ABC):
    @staticmethod
    @abstractmethod
    def compute_drift_evaluation(
        data_drift_context: DriftEvaluatorContext,
    ) -> DriftEvaluation:
        pass

class DriftEvaluation(TypedDict):
    should_alert: bool
    message: str
```

### Default Drift Evaluator

The default drift evaluator will return `should_alert = False`

### Alert Drift Evaluator

The Alert drift evaluator will reuturn `should_alert = True` for all drifts and a message containing the summary of the drift, example:

```
Drift detected:
- 🆕 0 addition
- ♻️ 2 modifications
- 🗑️ 0 deletion
```

To use the AlertDriftEvaluator, add it when you call snapshot_table like this:

```python
from driftdb.drift_evaluator.drift_evaluators import AlertDriftEvaluator

connector.snapshot_table(table_dataframe, table_name, drift_evaluator=AlertDriftEvaluator)

```

### Custom Drift Evaluator

You can provide a custom evaluator which is a function with a DriftEvaluatorContext containing the following properties:

```python
class DriftEvaluatorContext(TypedDict):
    before: pd.DataFrame
    after: pd.DataFrame
    summary: DriftSummary

class DriftSummary(TypedDict):
    added_rows: pd.DataFrame
    deleted_rows: pd.DataFrame
    modified_rows_unique_keys: pd.Index
    modified_patterns: pd.DataFrame

```

Then implement your class, and use it in snapshot_table.

```python
class MyDriftEvaluator(BaseDriftEvaluator):
    @staticmethod
    def compute_drift_evaluation(
        data_drift_context: DriftEvaluatorContext,
    ) -> DriftEvaluation:
        # do what you want
        if there_is_something_I_dont_like:
          return {"should_alert": True, "message": "No this should not happen"}
        return {"should_alert": False, "message": ""}
```

# CLI

Instead of storing data on github, you can store data locally and explore it with the cli.

# Getting started

## From dbt snapshot (dbt >= 1.6)

```shell
pip install driftdb

driftdb dbt snapshot
driftdb start
```

## From generated seeds

```shell
pip install driftdb

driftdb seed create
driftdb seed update

driftdb start
```

# Features

## Metrics

### Load a csv

```
driftdb load-csv path/to/csv
```

## Data visualization

```
driftdb start
```

Start the driftdb, and navigate to [localhost:9741/tables](http://localhost:9741/tables).
Visualize how a metric evolved, given a period, in a waterfall chart.
