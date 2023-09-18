import styled from "@emotion/styled";
import { useState } from "react";

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
  const [isChecked, setIsChecked] = useState(false);
  const [childChecked, setChildChecked] = useState(parentData.map(() => false));

  const handleChildChecked = (index: number) => {
    const newChildChecked = [...childChecked];
    newChildChecked[index] = !newChildChecked[index];
    setChildChecked(newChildChecked);
  };

  return (
    <DriftCardContainer>
      <div>
        <b>Filepath:</b> {filepath}
      </div>
      <div>
        <b>Period:</b> {periodKey}
      </div>
      <div>
        <b>Drift Date:</b>
        <input
          type="checkbox"
          checked={isChecked}
          onChange={() => setIsChecked(!isChecked)}
        />
        {new Date(driftDate).toLocaleString()}
      </div>
      {parentData.length > 0 && (
        <div style={{ display: "flex", flexDirection: "column" }}>
          <b style={{ alignSelf: "baseline" }}>Parent Data:</b>
          <ul>
            {parentData.map((parent, index) => (
              <li key={parent}>
                <input
                  type="checkbox"
                  checked={childChecked[index]}
                  onChange={() => handleChildChecked(index)}
                />
                {parent}
              </li>
            ))}
          </ul>
        </div>
      )}
    </DriftCardContainer>
  );
};
