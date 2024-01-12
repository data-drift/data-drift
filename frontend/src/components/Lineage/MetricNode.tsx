import styled from "@emotion/styled";
import { useState } from "react";

const StyledMetricNode = styled.div`
  border: 3px solid black;
  padding: 16px;
  border-radius: 0; /* Brutalist design prefers more angular, less rounded shapes. */
  min-width: 200px;
  background-color: ${(props) => props.theme.colors.text};
  text-align: left;
  max-height: 250px;
  overflow-y: scroll;
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
  transition: transform 0.2s, box-shadow 0.2s;
  &:hover {
    cursor: pointer !important;
    transform: scale(1.05);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
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
  onClick?: () => void;
};

export type LineageEvent = {
  type: "New Data" | "Drift";
  eventDate: Date | null;
  subEvents?: SubEvent[];
  onClick?: () => void;
};

type MetricNodeProps = {
  metricName: string;
  events?: LineageEvent[];
  eventsLoading?: boolean;
};

export const MetricNode = ({
  metricName = "",
  events = [],
  eventsLoading,
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
      event.onClick && event.onClick();
    }
  };

  return (
    <StyledMetricNode>
      <strong>{metricName}</strong>
      {eventsLoading ? (
        <small>Loading</small>
      ) : events.length > 0 ? (
        <ul>
          {events.map((event, index) => (
            <li key={index}>
              <EventChip
                eventType={event.type}
                onClick={() => handleEventClick(event)}
              >
                {event.type}{" "}
                {event.eventDate ? event.eventDate.toLocaleTimeString() : ""}
              </EventChip>
              {expandedEvent === event.type && (
                <ul>
                  {event.subEvents?.map((subEvent) => (
                    <li key={subEvent.name}>
                      <SubEventChip
                        onClick={() => subEvent.onClick && subEvent.onClick()}
                      >
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
