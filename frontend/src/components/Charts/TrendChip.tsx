import styled from "@emotion/styled";

type Props = {
  trend: "up" | "down" | "neutral";
  absoluteValue: number;
};

const ChipContainer = styled.div<{ trend: Props["trend"] }>`
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
`;

const TrendChip = ({ trend, absoluteValue }: Props) => {
  console.log(trend, absoluteValue);
  return (
    <ChipContainer trend={trend}>
      <strong>{absoluteValue}%</strong>
    </ChipContainer>
  );
};

export default TrendChip;
