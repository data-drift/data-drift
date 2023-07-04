import pandas as pd
import streamlit as st
from typing import List
from datetime import datetime
import plotly.express as px
import plotly.graph_objects as go
import numpy as np
import base64


from snapshot_utils import get_metric_by_month, get_snapshot_as_of_date, determine_type


@st.cache_data()
def compute_all_metric_day_by_day(local_df: pd.DataFrame, date_range: List[str]):
    print("Computing metric day by day...")
    all_monthly_metrics = pd.DataFrame()
    progress_text = "Operation in progress. Please wait."
    my_bar = st.progress(0, text=progress_text)
    index = 0
    for day in date_range:
        index += 1
        day_date = datetime.strptime(day, "%Y-%m-%d")
        current_df = get_snapshot_as_of_date(local_df, day_date)
        if current_df.empty:
            continue
        metric_by_month = get_metric_by_month(
            current_df, "metric_value", "metric_date", "%Y-%m"
        )
        metric_by_month["computation_day"] = day_date.date()
        all_monthly_metrics = pd.concat(
            [all_monthly_metrics, metric_by_month], ignore_index=True
        )
        progress = index / date_range.__len__()
        my_bar.progress(progress, text=progress_text)

    my_bar.empty()

    return all_monthly_metrics


def is_numeric_column(df, column_name):
    try:
        pd.to_numeric(df[column_name])
        return True
    except ValueError:
        return False


def is_date_column(df, column_name):
    try:
        pd.to_datetime(df[column_name])
        return True
    except ValueError:
        return False


@st.cache_data()
def parse_csv_file(uploaded_file):
    if uploaded_file is None:
        return pd.DataFrame()
    print("Parsing uploaded file...")
    local_df = pd.read_csv(uploaded_file, low_memory=False)

    print(type(local_df["dbt_valid_from"].iloc[0]))
    datetime_format = determine_type(local_df["dbt_valid_from"].iloc[0])
    print(f"Date format: {datetime_format}")
    # format dates
    try:
        if datetime_format == "timestamp":
            local_df["is_current_version"] = local_df["dbt_valid_to"].isnull()
            local_df["dbt_valid_from"] = pd.to_datetime(
                local_df["dbt_valid_from"], unit="s"
            )
            now = pd.Timestamp.now().timestamp()
            local_df["dbt_valid_to"] = local_df["dbt_valid_to"].fillna(now)
            local_df["dbt_valid_to"] = pd.to_datetime(
                local_df["dbt_valid_to"], unit="s"
            )
        else:
            local_df["is_current_version"] = local_df["dbt_valid_to"].isnull()
            local_df["dbt_valid_from"] = pd.to_datetime(local_df["dbt_valid_from"])
            now = pd.Timestamp.now().strftime("%Y-%m-%dT%H:%M:%S.%f")
            local_df["dbt_valid_to"] = local_df["dbt_valid_to"].fillna(now)
            local_df["dbt_valid_to"] = pd.to_datetime(local_df["dbt_valid_to"])
    except ValueError as e:
        invalid_row_index = int(str(e).split(" ")[-1])
        print(f"Invalid row index: {invalid_row_index}")

    print(f"Number of rows: {local_df.shape[0]}")
    print(f"Number of columns: {local_df.shape[1]}")
    return local_df


@st.cache_data()
def get_metric_value_and_date(local_df, metric_column, date_column):
    is_numeric_column_valid = is_numeric_column(local_df, metric_column)
    is_date_column_valid = is_date_column(local_df, date_column)
    if is_numeric_column_valid and is_date_column_valid:
        local_df["metric_value"] = local_df[metric_column]
        local_df["metric_date"] = local_df[date_column]

        local_df = local_df[
            (local_df["metric_date"].notnull()) & (local_df["metric_value"].notnull())
        ]
    else:
        st.write("Please select a valid metric column and a date column")
    return local_df


def run():
    st.sidebar.write("# dbt Snapshot analysis")
    st.sidebar.write(
        "*Changes of a metric’s historical value harms data consumers’ trust and is very difficult to solve for data teams.*"
    )
    st.sidebar.write(
        "***Upload a dbt snapshot of a table on which you build a key metric and see how you compare!***"
    )
    st.sidebar.write(
        "⭐️ [Star us on Github](https://github.com/data-drift/data-drift) ⭐️"
    )
    st.write("# 👋 Welcome")
    st.write("### Upload a dbt snapshot to get started")
    uploaded_file = st.file_uploader("Select a csv file")
    df = pd.DataFrame()
    if uploaded_file is None:
        st.write("### Analyse the stability of your key metrics")
        st.write(
            "Changes of a metric’s historical value harms data consumers’ trust and is very difficult to solve for data teams."
        )
        st.write(
            "**Upload a dbt snapshot of a table on which you build a key metric and see how you compare!**"
        )
        file_ = open("dbt-snapshot-analysis.gif", "rb")
        contents = file_.read()
        data_url = base64.b64encode(contents).decode("utf-8")
        file_.close()
        st.markdown(
            f'<img src="data:image/gif;base64,{data_url}" alt="cat gif">',
            unsafe_allow_html=True,
        )
    else:
        df = parse_csv_file(uploaded_file)
        unique_id_column = st.selectbox("Select unique id column", df.columns)

    if not df.empty:
        unique_ids = df.groupby(unique_id_column).size()
        versions_per_id = (
            df.groupby(unique_id_column).size().reset_index(name="version per id")
        )
        versions_count = versions_per_id["version per id"].value_counts().sort_index()
        lifespandf = df.copy()
        lifespandf["lifespan"] = (
            lifespandf["dbt_valid_to"] - lifespandf["dbt_valid_from"]
        )
        print(lifespandf["lifespan"].describe())
        lifespandf["lifespan_numeric"] = pd.to_numeric(
            lifespandf["lifespan"], errors="coerce"
        )
        lifespandf["lifespan (days)"] = lifespandf["lifespan"].dt.days
        print(lifespandf["lifespan (days)"].describe())

        all_versions = lifespandf["lifespan (days)"]
        dead_versions = lifespandf[lifespandf["is_current_version"] == False][
            "lifespan (days)"
        ]

        lifespan_df_with_alive = pd.DataFrame(
            dict(
                series=np.concatenate(
                    (
                        ["lifespan all rows (days)"] * len(all_versions),
                        ["lifespan without active (days)"] * len(dead_versions),
                    )
                ),
                data=np.concatenate((all_versions, dead_versions)),
            )
        )
        min_date = df["dbt_valid_from"].min()
        max_date = df["dbt_valid_from"].max()
        date_range = pd.date_range(start=min_date, end=max_date)
        date_range_str = date_range.strftime("%Y-%m-%d").tolist()

        metric_distribution_df = df.copy()
        st.write("### Compute the metric you want to analyse")

        col1, col2, col3 = st.columns(3)

        with col1:
            is_metric_or_count = st.radio(
                "Select the aggregation method", ["Count", "Sum"]
            )

        with col2:
            date_column = st.selectbox(
                "Select the reference date", metric_distribution_df.columns
            )
            is_date_column_valid = is_date_column(metric_distribution_df, date_column)
            if not is_date_column_valid:
                st.warning(
                    "Please select a valid date column, it should be supported by pandas.to_datetime."
                )

        with col3:
            is_numeric_column_valid = False
            if is_metric_or_count == "Sum":
                metric_column = st.selectbox(
                    "Select the column to sum", metric_distribution_df.columns
                )
                is_numeric_column_valid = is_numeric_column(
                    metric_distribution_df, metric_column
                )
                if not is_numeric_column_valid:
                    st.warning(
                        "Please select a valid metric column, it should be a numeric column."
                    )
            else:
                metric_distribution_df["count"] = 1
                metric_column = "count"
                is_numeric_column_valid = True

        if is_numeric_column_valid and is_date_column_valid:
            metric_distribution_df_with_formated_date = get_metric_value_and_date(
                metric_distribution_df, metric_column, date_column
            )

            all_results = compute_all_metric_day_by_day(
                metric_distribution_df_with_formated_date, date_range_str
            )

            all_results = all_results[all_results["metric_value"] > 0]

            all_results["latest_value"] = (
                all_results.sort_values(
                    ["metric_date", "computation_day"], ascending=[True, False]
                )
                .groupby("metric_date")["metric_value"]
                .transform("first")
            )

            all_results["relative_value"] = (
                all_results["metric_value"] / all_results["latest_value"]
            )

            volatility_filtered_df = all_results[all_results["metric_value"] > 0]
            volatility_grouped_df = volatility_filtered_df.groupby("metric_date")[
                "relative_value"
            ].apply(lambda x: ((x - 1) ** 2).mean() ** 0.5)
            all_monthly_volatility_from_latest_value = pd.DataFrame(
                volatility_grouped_df
            ).rename(columns={"relative_value": "volatility"})

            print("all_monthly_volatility", all_monthly_volatility_from_latest_value)

            volatility_grouped_df = all_results.groupby("metric_date")

            st.set_option("deprecation.showPyplotGlobalUse", False)
            fig = go.Figure()

            for metric_name, group in volatility_grouped_df:
                fig.add_trace(
                    go.Scatter(
                        x=group["computation_day"],
                        y=group["relative_value"],
                        name=metric_name,
                    )
                )

            fig.update_layout(
                xaxis_title="Computation Day",
                yaxis_title="Normalized Value",
                showlegend=True,
            )

            st.divider()
            st.write("### 💎 Summary")
            summary1, summary2 = st.columns(2)
            with summary1:
                st.write("#### Row lifespan")
                meanLifespan = dead_versions.mean()
                st.metric("Mean row lifespan", f"{meanLifespan:.2f} days")
                st.info(
                    "The lifespan is the number of days between two updates of a given row",
                    icon="💡",
                )

            with summary2:
                st.write("#### Metric volatility")
                meanVolatility = all_monthly_volatility_from_latest_value.mean().mean()
                st.metric(
                    "Average volatility (between 0 and 1)", f"{meanVolatility:.3f}"
                )
                st.info(
                    "The volatility shows the delta between computed value at time T and latest value",
                    icon="💡",
                )

            st.divider()
            st.write("### 🔎 Deep dive")
            st.write(
                "Navigate through the tabs below to get more insights from your dbt snapshot."
            )
            volatilityTab, lifespanTab, versionTab, dataTab = st.tabs(
                ["📊 Volatility", "⏱ Lifespan", "🗂 Version", "💿 Raw data"]
            )
            with volatilityTab:
                st.write("#### Metric volatility")
                st.write("##### Monthly volatility")
                st.plotly_chart(
                    px.line(all_monthly_volatility_from_latest_value)
                    .update_yaxes(range=[0, 1], title_text="Volatility")
                    .update_layout(showlegend=False)
                )
                volatility_expander = st.expander("💡 Learn more about volatility")
                volatility_expander.write(
                    "It shows how much your metric drifts. We compute the root mean square (RMS) of the difference between the computed value at time T and the latest computed value. It's ike a standard deviation but compared to the latest value instead of the mean. The closer to 0, the less your metric has drifted, the better."
                )
                st.write("##### Metric cohorts")
                st.plotly_chart(fig)

                normalized_value_expender = st.expander(
                    "💡 Learn more about normalized value"
                )
                normalized_value_expender.write(
                    "It shows how long your metric takes to stabilize to its current value. The normalized value is the value at time T divided by the latest known value. The time series of the normalized value converges to 1, the faster it stays there, the better."
                )

            with lifespanTab:
                st.write("#### Row lifespan")
                lifespan_fig = px.histogram(
                    lifespan_df_with_alive,
                    color="series",
                    barmode="overlay",
                    nbins=50,
                )
                st.plotly_chart(lifespan_fig, use_container_width=True)
                volatility_expander = st.expander("💡 Learn more about lifespan")
                volatility_expander.write("The lifespan is the number of days between two updates of a given row. For Lifespan without active, we exclude row that have never been updated (still valid).")

            with versionTab:
                st.write("#### Version distribution")
                print(f"Number of unique ids: {unique_ids.shape[0]}")
                st.bar_chart(versions_count)
                volatility_expander = st.expander("💡 Learn more about version")
                volatility_expander.write("A version represents a row for a given unique id. When an id has never been updated it's 1, the current version is the only version.")

            with dataTab:
                st.write("#### Raw data")
                st.write(
                    "The underlying raw data used to generate metrics and insights."
                )
                st.write(df)


def main():
    run()


if __name__ == "__main__":
    main()
