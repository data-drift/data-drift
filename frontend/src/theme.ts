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
    fontWeight: {
      thin: 100,
      extraLight: 200,
      light: 300,
      regular: 400,
      medium: 500,
      semiBold: 600,
      bold: 700,
      extraBold: 800,
      black: 900,
    },
    lineHeight: {
      normal: "normal",
    },
  },
};

export type ThemeType = typeof theme;
