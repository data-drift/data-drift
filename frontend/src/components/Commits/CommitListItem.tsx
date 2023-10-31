import styled from "@emotion/styled";

type CommitListItemProps = {
  type: "Drift" | "New Data";
  date: Date | null;
  name: string;
  commitUrl: string;
  isParentData: boolean;
};

const Container = styled.div`
  border: 1px solid #ccc;
  padding: 16px;
  border-radius: 0;
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
`;

const StyledButton = styled.button`
  padding: 8px 16px;
  background-color: #333;
  color: #fff;
  border-radius: 0px;
  border: 2px solid #fff;
  font-family: monospace;
`;

const FlexRow = styled.div`
  display: flex;
  justify-content: space-between;
  width: 100%;
`;

const FlexCenter = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

const TypeSpan = styled.span<{ type: CommitListItemProps["type"] }>`
  background-color: ${({ type, theme }) =>
    type === "Drift" ? theme.colors.primary : theme.colors.background};
  color: ${({ type }) => (type === "Drift" ? "#000" : "#fff")};
  border-radius: 0;
  padding: 4px 8px;
  font-weight: bold;
  margin-right: 8px;
`;

const ParentDataSpan = styled.span`
  background-color: ${({ theme }) => theme.colors.secondary};
  border-radius: 0;
  padding: 4px 8px;
  font-weight: bold;
`;

const DateSpan = styled.span`
  color: #888;
`;

export const CommitListItem = ({
  type,
  date,
  name,
  commitUrl,
  isParentData,
}: CommitListItemProps) => {
  return (
    <Container>
      <FlexRow>
        <FlexCenter>
          <TypeSpan type={type}>{type}</TypeSpan>
          {isParentData && <ParentDataSpan>Parent Data</ParentDataSpan>}
        </FlexCenter>
        {date && <DateSpan>{date.toLocaleString()}</DateSpan>}
      </FlexRow>
      <div>
        <p
          dangerouslySetInnerHTML={{ __html: name.replace(/\n/g, "<br/>") }}
        ></p>
      </div>
      {commitUrl && (
        <a href={commitUrl} target="_blank" rel="noopener noreferrer">
          <StyledButton>View</StyledButton>
        </a>
      )}
    </Container>
  );
};
