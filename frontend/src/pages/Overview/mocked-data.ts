export const mockedDiffTable = {
  tableProps1: {
    diffType: "removed",
    data: Array.from({ length: 130 }).map((_, i) => ({
      isEmphasized: i % 5 === 4,
      data: Array.from({ length: 10 }).map((_, j) => ({
        isEmphasized: i % 5 === 4 && (j + 2 * i) % 6 === 2,
        value: `Old ${i}-${j}`,
      })),
    })),
    headers: Array.from({ length: 10 }).map((_, j) => `Header ${j}`),
  },
  tableProps2: {
    diffType: "added",
    data: Array.from({ length: 130 }).map((_, i) => ({
      isEmphasized: i % 5 === 4,
      data: Array.from({ length: 10 }).map((_, j) => ({
        isEmphasized: i % 5 === 4 && (j + 2 * i) % 6 === 2,
        value: `New ${i}-${j}`,
      })),
    })),
    headers: Array.from({ length: 10 }).map((_, j) => `Header ${j}`),
  },
} as const;
