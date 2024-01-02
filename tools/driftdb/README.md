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

# Alerting

## Drift

A drift is a modification of historical data. It can be a modification, addition or deletion in a table that is supposed to be "non-moving data".

## Drift Handler

A drift handler is a function that conforms the type `DriftHandler`:

```python

DriftHandler = Callable[[DriftEvaluatorContext], DriftEvaluation]

# With DriftEvaluatorContext and DriftEvaluation being

class DriftEvaluatorContext:
    def __init__(self, before: pd.DataFrame, after: pd.DataFrame, summary: DriftSummary):
        self.before = before
        self.after = after
        self.summary = summary

class DriftEvaluation(TypedDict):
    should_alert: bool
    message: str
```

### Default Drift Handler

The default drift evaluator never triggers any alert, it returns `should_alert = False`.

### Alert Drift Handler

The `alert_drift_handler` will trigger an alert if there is a drifts, and an alert message containing the summary of the drift, example:

```
Drift detected:
- 🆕 0 addition
- ♻️ 2 modifications
- 🗑️ 0 deletion
```

To use the `alert_drift_handler`, add it when you call snapshot_table like this:

```python
from driftdb.alerting import alert_drift_handler

connector.snapshot_table(table_dataframe, table_name, drift_handler=alert_drift_handler)
```

### Threshold Drift Handler

The Threshold Drift Handler is designed to monitor changes in numerical values. It triggers an alert when a numerical value is updated and the absolute difference, when divided by the old value, exceeds a specified threshold.

Here's how you can use the Threshold Drift Handler:

```python
from driftdb.alerting import TresholdDriftHandlerFactory

# Set your desired threshold
threshold = 0.1  # Alert if the change is over 10%
threshold_handler = TresholdDriftHandlerFactory(numerical_cols=['metric_column_name'], threshold=threshold)

connector.snapshot_table(table_dataframe, table_name, drift_handler=threshold_handler)
```

### Custom Drift Handler

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

Then implement your handler, and use it in snapshot_table.

```python
def my_drift_handler(
    data_drift_context: DriftEvaluatorContext,
) -> DriftEvaluation:
    # do what you want
    if there_is_something_I_dont_like:
      return {"should_alert": True, "message": "No this should not happen"}
    return {"should_alert": False, "message": ""}
```

## New Data

When there is a new batch of data in a table, e.g. the results of last week. This addition is considered new data. It should not be confused with a new row entry of historical data. For instance, if a new transaction with a paying_date from a month ago is registered today, it will be considered a drift, not a new data.

## New Data Handler

You can also add alerting when inserting new data in the table.

### Detect Outlier Handler

You can create a `detect_outlier_handler` with the `DetectOutlierHandlerFactory` that takes 2 arguments, the numerical columns and the category columns.
For numerical columns, if the new data is an outlier (using the [interquartil method](https://en.wikipedia.org/wiki/Interquartile_range#Outliers)) it will trigger an alert.
For category columns, if a new category is detected, it will trigger an alert.

To use the `detect_outlier_handler`, add it when you call snapshot_table like this:

```python
from driftdb.alerting import DetectOutlierHandlerFactory
new_data_handler = DetectOutlierHandlerFactory(numerical_cols=["age"], categorical_cols=[])

connector.snapshot_table(table_dataframe, table_name, new_data_handler=new_data_handler)
```

### Custom New Data Handler

You can provide a custom handler which is a function with a NewDataEvaluatorContext containing the following properties:

```python
class NewDataEvaluatorContext:
    def __init__(self, before: pd.DataFrame, after: pd.DataFrame, added_rows: pd.DataFrame):
        self.before = before
        self.after = after
        self.added_rows = added_rows

```

Then implement your handler, and use it in snapshot_table.

```python
def my_new_data_handler(
    new_data_context: NewDataEvaluatorContext,
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
