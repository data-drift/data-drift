import { theme } from "../../theme";

type CommitListItemProps = {
  type: "Drift" | "New Data";
  date: Date | null;
  name: string;
  commitUrl: string;
  isParentData: boolean;
};

export const CommitListItem = ({
  type,
  date,
  name,
  commitUrl,
  isParentData,
}: CommitListItemProps) => {
  return (
    <div
      style={{
        border: "1px solid #ccc",
        padding: "16px",
        borderRadius: "0",
        marginBottom: "16px",
        display: "flex",
        flexDirection: "column",
        alignItems: "flex-start",
      }}
    >
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          width: "100%",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: "4px" }}>
          <span
            style={{
              backgroundColor:
                type === "Drift"
                  ? theme.colors.primary
                  : theme.colors.background,
              color: type === "Drift" ? "#000" : "#fff",
              borderRadius: "0",
              padding: "4px 8px",
              fontWeight: "bold",
            }}
          >
            {type}
          </span>
          {isParentData && (
            <span
              style={{
                backgroundColor: theme.colors.secondary,
                borderRadius: "0",
                padding: "4px 8px",
                fontWeight: "bold",
              }}
            >
              Parent Data
            </span>
          )}
        </div>
        {date && <span style={{ color: "#888" }}>{date.toLocaleString()}</span>}
      </div>
      <div style={{ marginTop: "8px" }}>
        <p
          dangerouslySetInnerHTML={{ __html: name.replace(/\n/g, "<br/>") }}
        ></p>
      </div>
      <a href={commitUrl} target="_blank" rel="noopener noreferrer">
        <button
          style={{
            padding: "8px 16px",
            backgroundColor: "#333",
            color: "#fff",
            borderRadius: "0px",
            border: "2px solid #fff",
            fontFamily: "monospace",
          }}
        >
          View Commit
        </button>
      </a>
    </div>
  );
};
