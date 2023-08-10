</br>
<p align="center">
  <a href="https://www.data-drift.io">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="./datadrift-logo-light.png" width="200px">
      <source media="(prefers-color-scheme: light)" srcset="./datadrift-logo-dark.png" width="200px">
      <img src="./datadrift-logo-dark.png" width="200px" alt="Datadrift logo" />
    </picture>
  </a>
</p>

<p align="center">
  <a href="https://github.com/data-drift/data-drift/actions/workflows/datadrift-build.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/data-drift/data-drift/datadrift-build.yml?style=flat-square" alt="Data-Drift Build">
  </a>
  <a href="https://main--64be84b7fe2172aa386216b8.chromatic.com/?path=/story/drift-dualtable--simple-case">
    <img src="https://img.shields.io/badge/storybook-visit-FF4785.svg?style=flat-square&logo=storybook" alt="Storybook">
  </a>
  <a href="https://pypi.org/project/datagit/">
    <img src="https://img.shields.io/pypi/v/datagit?style=flat-square" alt="DataGit version">
  </a>
</p>

<h1 align="center" >The Context Layer for your Metrics</h1>
<p align="center">Supercharge your semantic layer with context for your data consumers.
Build actionnable and trusted metrics with changelog, monitoring and centralized interpretation.
</p>

<p align="center"><a href="https://data-drift.io">Website</a> · <a href="https://www.data-drift.io/blog">Blog</a> · <a href="https://github.com/data-drift/data-drift/issues">Issues</p>
</br>

<p align="center">
  <a href="https://www.data-drift.io">
    <img src="./datadrift-new-drift.png" alt="DataDrift hero with metric volatility charts" />
  </a>
</p>
</br>

# 👋 About

## The Problem: Data consumers interact with raw metrics only. Context is key to make metrics trusted and actionable.

</br>
Most data consumers never query a data warehouse table, yet use data on a daily basis through the lens of metrics.

To trust metrics and make decisions based on them, we need to guarante the quality of the metrics itself (not only the underlying tables) and give (a lot of) context around it.

</br>
Context is vital because is gives data consumers awareness of:

- **Computation**: how a metric was calculated and associated caveats

- **Governance**: who computed the metric, who reviewed it and who acts on it

- **Changelog**: when was the metric computed, when was it last updated and what was the impact of the change

- **Historical trajectory**: what happened to the metric overtime to prevent misinterpretation (remember the reason of that single unexpected revenue drop last quarter?)

</br>
<p align="center">
  <a href="https://www.data-drift.io">
    <img src="./datadrift-repo-meme.png" alt="DataDrift hero with metric volatility charts" />
  </a>
</p>
Data teams be like...(yes, providing context is hard)

</br>

## The Solution: Datadrift, the open-source context layer for data-driven companies

- Comprehensive **metadata** to certify a quality standard for metrics

- Easy governance with **version control and historisation**

- Usable **changelog and audit trail** of a metric lifecycle for data consumers

- Visibility and **centralized knowledge** of a metric’s historical trajectory

</br>

**Open-source, Open Architecture:**

- **Flexible**: compose your own context layer based on our building blocks

- **Secure**: deploy on your own infra to keep control over your data

- **Integrated**: not another tool to manage in your stack, use datadrift directly from current tools (dbt, BI)

</br>
<p align="center">
  <a href="https://www.data-drift.io">
    <img src="./datadrift-stack-schema.png" alt="Headless context for your metrics, wherever they are" />
  </a>
</p>

</br>

# 🚀 Quickstart

## Version-control your key metrics with Datagit

[Install Datagit](https://github.com/data-drift/data-drift/tree/main/tools/datagit#datagit) to historise and diff-checks your metrics' underlying data.

This is a mandatory step to generate context for your metrics. You can [learn more about Datagit in this article](https://www.data-drift.io/blog/git-for-your-data).

## Deploy Datadrift locally

Follow our [step-by-step installation guide](https://lucas2vries.notion.site/Step-by-Step-Installation-752ffb590d4e4b27bdb753f9654ef676) to use Datadrift.

## Use our cloud-based product

[Contact our team by filling the form on our website](https://www.data-drift.io/) to get started with Datadrift Cloud.

</br>

# 💚 Helping us

We 💚 contributions big and small. In priority order (although everything is appreciated) with the most helpful first:

- Star this repo to help us get visibility
- [Become a Design Partner](https://www.data-drift.io/design-partner) to co-built a product you & users love
- [Open an issue](https://github.com/data-drift/data-drift/issues) to share your idea or a bug you might have spotted

</br>

# 🗓 Upcoming features

## Coming Soon

🌀 Automatic lineage drill-down and diff checks. [Learn more about this feature](https://www.data-drift.io/join-the-waitlist)

🌀 dbt integration

## Coming later this year

🗓 Sharing via Slack & emails

🗓 Warehouse native integration

🗓 BI tools integration

🗓 Gsheet integration

Track planning on [Github Projects](https://github.com/orgs/data-drift/projects/3) and help us prioritising by upvoting or creating [issues](https://github.com/data-drift/data-drift/issues).
