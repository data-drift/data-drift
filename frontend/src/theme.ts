const diffColors = {
  strongPositive: "#008F4B",
  lightPositive: "#124b30",
  strongNegative: "#8A1515",
  lightNegative: "#491a1b",
};

const hexToRgba = (hex: string, alpha = 1): string => {
  // Ensure the hex value is 6 characters long, ignoring possible leading '#'
  const validHex = hex.slice(-6);

  // Extract the red, green, and blue components
  const [r, g, b] = validHex.match(/\w\w/g)!.map((val) => parseInt(val, 16));

  return `rgba(${r},${g},${b},${alpha})`;
};

export const theme = {
  colors: {
    background: "#1E1E1E",
    background2: "#2C2C2C",
    text: "#F0F0F0",
    text2: "#FFFFFF",
    ...diffColors,
    hexToRgba,
    diffSelect(diffType: "added" | "removed") {
      return diffType === "added"
        ? { strong: diffColors.strongPositive, light: diffColors.lightPositive }
        : {
            strong: diffColors.strongNegative,
            light: diffColors.lightNegative,
          };
    },
    charts: {
      2022: ["#62A8E9", "#6EDB9F"],
      2023: ["#6EDB9F", "#287D57"],
    },
    dataUp: "#006DCC",
    dataDown: "#E35300",
    primary: "#8FFF00",
    secondary: "#FF00FF",
  },
  spacing: (multiplicator: number) => `${multiplicator * 4}px`,
  text: {
    fontSize: {
      small: "10px",
      medium: "16px",
    },
  },
  upLeftClipping: "polygon(12px 0, 100% 0, 100% 100%, 0 100%, 0 12px)",
  downLeftClipping:
    "polygon(0 0, 100% 0, 100% 100%, 12px 100%, 0 calc(100% - 12px))",
};

export type ThemeType = typeof theme;
