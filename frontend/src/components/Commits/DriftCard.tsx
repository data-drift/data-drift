import styled from "@emotion/styled";

const DriftCardContainer = styled.div`
  border: 1px solid #ccc;
  border-radius: 0px;
  padding: 8px;
  margin-bottom: 16px;
  height: fit-content;
  display: flex;
  flex-direction: column;
  align-items: start;
`;

export const DriftCard = ({
  filepath,
  periodKey,
  driftDate,
  parentData,
}: {
  filepath: string;
  periodKey: string;
  driftDate: string;
  parentData: string[];
}) => {
  return (
    <DriftCardContainer>
      <div>
        <b>Filepath:</b> {filepath}
      </div>
      <div>
        <b>Period:</b> {periodKey}
      </div>
      <div>
        <b>Drift Date:</b> {new Date(driftDate).toLocaleString()}
      </div>
      {parentData.length > 0 && (
        <div>
          <b>Parent Data:</b>
          <ul>
            {parentData.map((parent) => (
              <li key={parent}>{parent}</li>
            ))}
          </ul>
        </div>
      )}
    </DriftCardContainer>
  );
};
