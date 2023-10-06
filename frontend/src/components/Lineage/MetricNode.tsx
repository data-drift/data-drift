import styled from "@emotion/styled";
import { useState } from "react";

const StyledMetricNode = styled.div`
  border: 3px solid black;
  padding: 16px;
  border-radius: 0; /* Brutalist design prefers more angular, less rounded shapes. */
  min-width: 200px;
  background-color: ${(props) => props.theme.colors.text};
  text-align: left;
  strong {
    display: block;
    margin-bottom: 8px;
    font-size: 20px;
    color: ${(props) => props.theme.colors.background};
    text-transform: uppercase;
  }
  small {
    color: ${(props) => props.theme.colors.background2};
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

const EventChip = styled.div<{ eventType: LineageEvent["type"] }>`
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

const SubEventChip = styled.div`
  background-color: #f0f0f0; // Lighter shade for differentiation
  padding: 4px 8px;
  border-radius: 0;
  margin-top: 5px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #e0e0e0; // Slight change on hover for interactivity
  }
`;

type SubEvent = {
  name: string;
};

export type LineageEvent = {
  type: "New Data" | "Drift";
  subEvents?: SubEvent[];
};

type MetricNodeProps = {
  metricName: string;
  events?: LineageEvent[];
};

export const MetricNode = ({
  metricName = "",
  events = [],
}: MetricNodeProps) => {
  const [expandedEvent, setExpandedEvent] = useState<
    LineageEvent["type"] | null
  >(null);

  const handleEventClick = (event: LineageEvent) => {
    const hasSubEvents = event.subEvents?.length && event.subEvents?.length > 0;
    if (hasSubEvents) {
      if (expandedEvent === event.type) {
        setExpandedEvent(null);
      } else {
        setExpandedEvent(event.type);
      }
    } else {
      alert(event);
    }
  };

  return (
    <StyledMetricNode>
      <strong>{metricName}</strong>
      {events.length > 0 ? (
        <ul>
          {events.map((event, index) => (
            <li key={index}>
              <EventChip
                eventType={event.type}
                onClick={() => handleEventClick(event)}
              >
                {event.type}
              </EventChip>
              {expandedEvent === event.type && (
                <ul>
                  {event.subEvents?.map((subEvent) => (
                    <li key={subEvent.name}>
                      <SubEventChip onClick={() => alert(subEvent.name)}>
                        {subEvent.name}
                      </SubEventChip>
                    </li>
                  ))}
                </ul>
              )}
            </li>
          ))}
        </ul>
      ) : (
        <small>No Events</small>
      )}
    </StyledMetricNode>
  );
};

export default MetricNode;
