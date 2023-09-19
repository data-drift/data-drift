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

const useDriftCard = ({
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
  return {
    filepath,
    periodKey,
    driftDate: { driftDate, isChecked, setIsChecked },
    parentData: {
      parentData,
      childChecked,
      setChildChecked,
      handleChildChecked,
    },
  };
};

type DriftCardState = ReturnType<typeof useDriftCard>;

const DriftCard = (state: DriftCardState) => {
  return (
    <DriftCardContainer>
      <div>
        <b>Filepath:</b> {state.filepath}
      </div>
      <div>
        <b>Period:</b> {state.periodKey}
      </div>
      <div>
        <b>Drift Date:</b>
        <input
          type="checkbox"
          checked={state.driftDate.isChecked}
          onChange={() =>
            state.driftDate.setIsChecked(!state.driftDate.isChecked)
          }
        />
        {new Date(state.driftDate.driftDate).toLocaleString()}
      </div>
      {state.parentData.parentData.length > 0 && (
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            textAlign: "left",
          }}
        >
          <b style={{ alignSelf: "baseline" }}>Parent Data:</b>
          <ul>
            {state.parentData.parentData.map((parent, index) => (
              <li key={parent}>
                <input
                  type="checkbox"
                  checked={state.parentData.childChecked[index]}
                  onChange={() => state.parentData.handleChildChecked(index)}
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

DriftCard.useState = useDriftCard;

export default DriftCard;
