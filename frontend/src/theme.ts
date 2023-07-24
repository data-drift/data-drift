export const theme = {
  colors: {
    background: "#1E1E1E",
    background2: "#2C2C2C",
    text: "#F0F0F0",
    text2: "#FFFFFF",
  },
  spacing: (multiplicator: number) => `${multiplicator * 4}px`,
  text: {
    fontSize: {
      small: "10px",
      medium: "16px",
    },
  },
};

export type ThemeType = typeof theme;
