import styled from "@emotion/styled";

type Props = {
  trend: "up" | "down" | "neutral";
  absoluteValue: number;
};

const ChipContainer = styled.div<{ trend: Props["trend"] }>`
  display: inline-block;
  background-color: ${({ theme, trend }) =>
    trend === "up"
      ? theme.colors.dataUp
      : trend === "down"
      ? theme.colors.dataDown
      : theme.colors.text};
  padding: ${({ theme }) => theme.spacing(1)} ${({ theme }) => theme.spacing(2)};
  color: ${({ theme, trend }) =>
    trend === "neutral" ? "black" : theme.colors.text};
  clip-path: ${({ theme, trend }) =>
    trend === "up" ? theme.upLeftClipping : theme.downLeftClipping};

  &::after {
    position: relative;
    content: "${({ trend }) =>
      trend === "up" ? "↗" : trend === "down" ? "↘" : ""}";
    bottom: ${({ trend }) =>
      trend === "up" ? "9px" : trend === "down" ? "-5px" : ""};
    right: -5px;
    width: 0;
    height: 0;
    font-size: 10px;
  }
`;

const TrendChip = ({ trend, absoluteValue }: Props) => {
  return (
    <ChipContainer trend={trend}>
      <strong>{absoluteValue}%</strong>
    </ChipContainer>
  );
};

export default TrendChip;
