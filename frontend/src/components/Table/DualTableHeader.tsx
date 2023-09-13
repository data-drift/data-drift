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

export const DualTableHeader = () => {
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");

  const handleStartDateChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    setStartDate(e.target.value);
  };

  const handleEndDateChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    setEndDate(e.target.value);
  };
  return (
    <StyledDatePicker>
      <StyledDateInput
        type="date"
        value={startDate}
        onChange={handleStartDateChange}
      />
      <StyledDivider>to</StyledDivider>
      <StyledDateInput
        type="date"
        value={endDate}
        onChange={handleEndDateChange}
      />
    </StyledDatePicker>
  );
};
