import { extractFileNameAndPath } from "./string-helpers";

describe("extractFileNameAndPath", () => {
  it("should extract the file name and path from a file path string", () => {
    const filepath = "/path/to/myfile.txt";
    const result = extractFileNameAndPath(filepath);
    expect(result).toEqual({
      path: "/path/to",
      fileName: "myfile",
      extension: "txt",
    });
  });

  it("should handle purefile names", () => {
    const filepath = "myfile.txt";
    const result = extractFileNameAndPath(filepath);
    expect(result).toEqual({
      path: "",
      fileName: "myfile",
      extension: "txt",
    });
  });

  it("should handle path without extension", () => {
    const filepath = "/path/to/myfile";
    const result = extractFileNameAndPath(filepath);
    expect(result).toEqual({
      path: "/path/to",
      fileName: "myfile",
      extension: "",
    });
  });
});
