# DriftDB

<p align="center">
  <a href="https://pypi.org/project/driftdb/">
    <img src="https://img.shields.io/pypi/v/driftdb?style=flat-square" alt="driftdb version">
  </a>
  <a href="https://pypi.org/project/driftdb/">
    <img src="https://img.shields.io/pypi/dm/driftdb?style=flat-square" alt="driftdb monthly downloads">
  </a>
</p>

# Purpose

Driftdb is an historical metric store. Instead of providing evolution of a metric through time, it provides the
evolution of a **measurement**.
e.g `mrr of june` was 100, now it is 99.

It's designed to address stability of a metric, which is directly correlated to trustability.

# Getting started

## From dbt snapshot (dbt >= 1.6)

```shell
pip install driftdb

driftdb dbt snapshot
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

Start the driftdb, and navigate to [localhost:9741/tables](http://localhost:9741/tables).
Visualize how a metric evolved, given a period, in a waterfall chart.
