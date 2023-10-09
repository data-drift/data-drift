import styled from "@emotion/styled";

export const Container = styled.div`
  width: 100%;
  box-sizing: border-box;
`;

export const StyledHeader = styled.header`
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: ${(props) =>
    props.theme.colors.background2}; // muted background color
  padding: 20px 40px;
  width: 100%;
  box-sizing: border-box;
`;

export const StyledDate = styled.div`
  font-size: 32px;
  font-weight: bold;
  color: ${(props) => props.theme.colors.text};
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 4px;
`;

export const StyledSelect = styled.select`
  background-color: ${(props) => props.theme.colors.background};
  border: none;
  padding: 10px 20px;
  font-size: 16px;
  cursor: pointer;
`;

export const StyledDateButton = styled.button`
  background-color: transparent;
  color: ${(props) => props.theme.colors.text2};
  border: 0;
  border-radius: 0;
  padding: 5px 10px;
  font-size: 18px;
  cursor: pointer;
  &:hover {
    background-color: #bbb; // or any other color indication for hover
  }
`;

export const LineageContainer = styled.div`
  background-color: ${(props) => props.theme.colors.background2};
  text-align: left;
`;

export const StyledCollapsibleTitle = styled.button`
  cursor: pointer;
  padding: 10px;
  border: none;
  background-color: ${(props) => props.theme.colors.background2};
`;

export const StyledCollapsibleContent = styled.div<{ isCollapsed: boolean }>`
  height: ${(props) => (props.isCollapsed ? "0" : "300px")};

  overflow: hidden;
`;

export const DiffTableContainer = styled.div``;
