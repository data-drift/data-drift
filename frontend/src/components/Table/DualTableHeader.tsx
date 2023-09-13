import { ChangeEventHandler, useState } from "react";
import styled from "@emotion/styled";

const StyledDivider = styled.span`
  color: #fff;
`;

const StyledDatePicker = styled.div`
  display: flex;
  align-items: center;
  align-self: flex-start;
`;

const StyledDateInput = styled.input`
  background-color: ${(props) => props.theme.colors.background2};
  border: 1px solid ${(props) => props.theme.colors.text};
  color: #fff;
  padding: 8px;
  margin: 8px;
`;

const useDualTableHeader = () => {
  const searchParams = new URLSearchParams(window.location.search);
  const urlStartDate = searchParams.get("startDate");
  const urlEndDate = searchParams.get("endDate");
  const [startDate, setStartDate] = useState(urlStartDate || "");
  const [endDate, setEndDate] = useState(urlEndDate || "");

  const handleStartDateChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    const newStartDate = e.target.value;
    const searchParams = new URLSearchParams(window.location.search);
    searchParams.set("startDate", newStartDate);
    const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
    window.history.pushState({ path: newUrl }, "", newUrl);
    setStartDate(newStartDate);
  };

  const handleEndDateChange: ChangeEventHandler<HTMLInputElement> = (e) => {
    const newEndDate = e.target.value;
    const searchParams = new URLSearchParams(window.location.search);
    searchParams.set("endDate", newEndDate);
    const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
    console.log(newUrl);
    window.history.pushState({ path: newUrl }, "", newUrl);
    setEndDate(newEndDate);
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
