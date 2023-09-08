type CommitListItemProps = {
  type: "Drift" | "New Data";
  date: Date;
  name: string;
  commitUrl: string;
};

export const CommitListItem = ({
  type,
  date,
  name,
  commitUrl,
}: CommitListItemProps) => {
  return (
    <div
      style={{
        border: "1px solid #ccc",
        padding: "16px",
        borderRadius: "0",
        marginBottom: "16px",
      }}
    >
      <div style={{ display: "flex", justifyContent: "space-between" }}>
        <div style={{ display: "flex", alignItems: "center", gap: "4px" }}>
          <span
            style={{
              backgroundColor: type === "Drift" ? "#e0e0e0" : "#a5d6a7",
              color: type === "Drift" ? "#000" : "#fff",
              borderRadius: "0",
              padding: "4px 8px",
              fontWeight: "bold",
            }}
          >
            {type}
          </span>
        </div>
        <span style={{ color: "#888" }}>{date.toLocaleString()}</span>
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
