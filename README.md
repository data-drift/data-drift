# Datadrift is an open-source metric-focused data quality tool
- Historize your key metrics with kpi-git-history to get started quickly
- Monitor unexpected raw data updates impacting an historical metric  
- Investigate simply what data has changed and how the metric has been impacted  
- Report automatically why the metric has drifted via a shared Notion page with data consumers 
- Keep control over your data by deploying on your infrastructure 

# Get started for free
We are just launching our beta and are looking for feedbacks.

# Contributing
We <3 contributions big and small. In priority order (although everything is appreciated) with the most helpful first:

- Give us feedback by filling this [form](https://forms.gle/8q2NzoZC417cgdo38)
- Submit a feature request or bug report directly in this [form](https://forms.gle/8q2NzoZC417cgdo38)

# Philosophy 
Accurate warehouse data is essential but current tools only address row-level quality (e.g. non-null values), not time-varying metrics like ARR or revenues. These metrics inform critical decisions for internal teams, investors, and public markets. Financial metrics in particular demand precision and immutability.

However, metrics from your data warehouse are mutable due to:
- human errors (engineering team has bugs, third party API has bugs);
- breaking changes (definition of KPI is changed); and
- delay to reach its final state (eg. a transaction status assumed to be completed but canceled several weeks after).

As a result, data consumers (finance teams, accounting, executives or investors) consistently require data analysts to investigate “data quality” issues. The data analyst becomes an auditor.

Our ultimate goal: ***transform the data warehouse into a reliable, audit-ready system for crucial decision-making, with 100% accurate and trusted metrics.***
