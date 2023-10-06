import styled from "@emotion/styled";

const StyledMetricNode = styled.div`
  border: 3px solid black;
  padding: 16px;
  border-radius: 0; /* Brutalist design prefers more angular, less rounded shapes. */
  min-width: 200px;
  background-color: ${(props) => props.theme.colors.text};
  strong {
    display: block;
    margin-bottom: 8px;
    font-size: 20px;
    color: ${(props) => props.theme.colors.background};
    text-transform: uppercase;
  }
  ul {
    list-style-type: square; /* A more rugged list bullet. */
    padding-left: 20px;
  }

  li {
    margin-bottom: 4px;
    color: ${(props) => props.theme.colors.background2};
  }
`;

type Event = {
  type: "New Data" | "Drift";
};

type MetricNodeProps = {
  metricName: string;
  events?: Event[];
};

export const MetricNode = ({
  metricName = "",
  events = [],
}: MetricNodeProps) => {
  return (
    <StyledMetricNode>
      <strong>{metricName}</strong>
      <ul>
        {events.map((event, index) => (
          <li key={index}>{event.type}</li>
        ))}
      </ul>
    </StyledMetricNode>
  );
};

export default MetricNode;
