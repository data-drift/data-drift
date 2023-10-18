import { ChangeEventHandler, useEffect, useState } from "react";
import { Params, useLoaderData } from "react-router-dom";
import {
  WaterfallChart,
  WaterfallChartProps,
} from "../components/Charts/WaterfallChart";
import Loader from "../components/Common/Loader";
import { getMetricHistory } from "../services/data-drift";
import { getNiceTickValues } from "recharts-scale";
import { theme } from "../theme";

const loader = ({ params }: { params: Params<"metricName" | "tableName"> }) => {
  return { params };
};

type LoaderData = Awaited<ReturnType<typeof loader>>;

type Mutable<T> = {
  -readonly [P in keyof T]: T[P];
};

const getWaterfallChartPropsFromMetricHistory = (
  metricMeasurements: {
    LineCount: number;
    Metric: string;
    MeasurementMetaData: {
      MeasurementTimestamp: number;
      MeasurementDate: string;
      MeasurementDateTime: string;
      MeasurementComments: {
        CommentAuthor: string;
        CommentBody: string;
      }[];

      IsMeasureAfterPeriod: boolean;
      MeasurementId: string;
    };
  }[]
): WaterfallChartProps => {
  metricMeasurements.sort((aMeasurement, bMeasurement) => {
    return (
      aMeasurement.MeasurementMetaData.MeasurementTimestamp -
      bMeasurement.MeasurementMetaData.MeasurementTimestamp
    );
  });

  const data = [] as Mutable<WaterfallChartProps["data"]>;

  metricMeasurements.forEach((measurement, index) => {
    const commitDate = new Date(
      measurement.MeasurementMetaData.MeasurementDateTime
    );
    const formatedDate = `${String(commitDate.getMonth() + 1).padStart(
      2,
      "0"
    )}-${String(commitDate.getDate()).padStart(2, "0")}`;
    if (index === 0) {
      const yMin = Math.min(
        ...metricMeasurements.map((measurement) =>
          parseFloat(measurement.Metric)
        )
      );
      const yMax = Math.max(
        ...metricMeasurements.map((measurement) =>
          parseFloat(measurement.Metric)
        )
      );
      const niceTicks = getNiceTickValues([yMin, yMax], 5);
      data.push({
        day: formatedDate,
        drift: [niceTicks[0], parseFloat(measurement.Metric)],
        fill: theme.colors.text,
        isInitial: true,
      });
    } else {
      const latestValue = parseFloat(metricMeasurements[index - 1].Metric);
      const newValue = parseFloat(measurement.Metric);
      if (latestValue === newValue) {
        return;
      }
      data.push({
        day: formatedDate,
        drift: [latestValue, newValue],
        fill:
          latestValue > newValue ? theme.colors.dataDown : theme.colors.dataUp,
      });
    }
  });

  return { data };
};

const MetricPage = () => {
  const loader = useLoaderData() as LoaderData;

  const searchParams = new URLSearchParams(window.location.search);
  const initialPeriodKey = searchParams.get("periodKey");
  const [periodKey, setPeriodKey] = useState(initialPeriodKey || "2023-01");

  const handlePeriodKeyChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    const newPeriodKey = e.target.value;
    const searchParams = new URLSearchParams(window.location.search);
    searchParams.set("periodKey", newPeriodKey);
    const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
    window.history.pushState({ path: newUrl }, "", newUrl);
    setPeriodKey(newPeriodKey);
  };

  const [metricHistory, setMetricHistory] = useState<{
    loading: boolean;
    waterfallChartProps: WaterfallChartProps;
  }>({ loading: false, waterfallChartProps: { data: [] } });

  useEffect(() => {
    setMetricHistory({ loading: true, waterfallChartProps: { data: [] } });
    void getMetricHistory({
      metric: loader.params.metricName || "",
      periodKey,
      store: "default",
      table: loader.params.tableName || "",
    }).then(
      (result) => {
        const waterfallChartProps = getWaterfallChartPropsFromMetricHistory(
          result.data.metricHistory
        );
        setMetricHistory({
          loading: false,
          waterfallChartProps,
        });
      },
      (error) => console.log(error)
    );
  }, [periodKey, loader.params]);

  return (
    <div style={{ height: "100vh" }}>
      <h1>Metric: {loader.params.metricName}</h1>
      <label htmlFor="periodKey">Period: </label>
      <input
        type="text"
        id="periodKey"
        name="periodKey"
        value={periodKey}
        onChange={handlePeriodKeyChange}
        title="Period Key"
      />
      <span style={{ marginLeft: "0.5rem", color: "gray" }}>
        (accepted: YYYY-MM-DD YYYY-MM YYYY-Q1 YYYY-W01 YYYY)
      </span>
      <div style={{ paddingTop: "16px" }}>
        {metricHistory?.loading ? (
          <Loader />
        ) : (
          metricHistory && (
            <WaterfallChart {...metricHistory.waterfallChartProps} />
          )
        )}
      </div>
    </div>
  );
};

MetricPage.loader = loader;

export default MetricPage;
