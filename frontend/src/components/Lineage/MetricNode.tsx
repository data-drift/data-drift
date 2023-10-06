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

const EventChip = styled.div<{ eventType: Event["type"] }>`
  background-color: ${(props) =>
    props.eventType === "Drift"
      ? props.theme.colors.primary
      : props.theme.colors.background};
  color: ${(props) => (props.eventType === "Drift" ? "#000" : "#fff")};
  border-radius: 0;
  padding: 4px 8px;
  display: inline-block;
  text-align: center;
  margin-right: 5px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s; // For a smooth visual feedback
  &:hover {
    transform: scale(1.05); // Slightly enlarges the chip
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1); // Adds a subtle shadow for depth
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
          <li key={index}>
            <EventChip eventType={event.type}>{event.type}</EventChip>
          </li>
        ))}
      </ul>
    </StyledMetricNode>
  );
};

export default MetricNode;
