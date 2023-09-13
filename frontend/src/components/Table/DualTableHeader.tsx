import { ChangeEventHandler, useState } from "react";
import styled from "@emotion/styled";

const StyledDivider = styled.span`
  color: #fff;
`;

const StyledDatePicker = styled.div`
  display: flex;
  align-items: center;
`;

const StyledDateInput = styled.input`
  background-color: #444;
  border: 2px solid #fff;
  color: #fff;
  padding: 8px;
  margin: 8px;
`;

const useDualTableHeader = () => {
  const [startDate, setStartDate] = useState("2023-06-01");
  const [endDate, setEndDate] = useState("2023-07-01");

  const handleStartDateChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    setStartDate(e.target.value);
  };

  const handleEndDateChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    setEndDate(e.target.value);
  };
  return {
    startDate,
    endDate,
    handleStartDateChange,
    handleEndDateChange,
  };
};

type DualTableHeaderState = ReturnType<typeof useDualTableHeader>;

const DualTableHeader = ({ state }: { state: DualTableHeaderState }) => {
  const { startDate, endDate, handleStartDateChange, handleEndDateChange } =
    state;

  return (
    <StyledDatePicker>
      <StyledDateInput
        type="date"
        value={startDate}
        onChange={handleStartDateChange}
        title="Start date included"
      />
      <StyledDivider>to</StyledDivider>
      <StyledDateInput
        type="date"
        value={endDate}
        onChange={handleEndDateChange}
        title="End date excluded"
      />
    </StyledDatePicker>
  );
};

DualTableHeader.useState = useDualTableHeader;

export default DualTableHeader;
