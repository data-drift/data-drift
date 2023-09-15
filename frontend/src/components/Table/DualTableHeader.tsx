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

const StyledUderlineButton = styled.span`
  text-decoration: underline;
  color: ${(props) => props.theme.colors.text};
  cursor: pointer;
  margin-left: 16px;
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

  const clearFilters = () => {
    const searchParams = new URLSearchParams(window.location.search);
    searchParams.delete("startDate");
    searchParams.delete("endDate");
    const newUrl = `${window.location.pathname}?${searchParams.toString()}`;
    window.history.pushState({ path: newUrl }, "", newUrl);
    setStartDate("");
    setEndDate("");
  };

  return {
    startDate,
    endDate,
    handleStartDateChange,
    handleEndDateChange,
    clearFilters,
  };
};

type DualTableHeaderState = ReturnType<typeof useDualTableHeader>;

const DualTableHeader = ({
  state,
  copyAction,
}: {
  state: DualTableHeaderState;
  copyAction: () => void;
}) => {
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
      {(startDate != "" || endDate != "") && (
        <StyledUderlineButton onClick={state.clearFilters}>
          Clear Filters
        </StyledUderlineButton>
      )}
      <StyledUderlineButton onClick={() => copyAction()}>
        Copy Table to Clipboard
      </StyledUderlineButton>
    </StyledDatePicker>
  );
};

DualTableHeader.useState = useDualTableHeader;

export default DualTableHeader;
