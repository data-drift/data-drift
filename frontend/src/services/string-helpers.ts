export const extractFileNameAndPath = (filepath: string) => {
  const lastSlashIndex = filepath.lastIndexOf("/");
  const path = filepath.substring(0, lastSlashIndex);
  const fileNameWithExtension = filepath.substring(lastSlashIndex + 1);
  const dotIndex = fileNameWithExtension.lastIndexOf(".");
  const fileName =
    dotIndex >= 0
      ? fileNameWithExtension.substring(0, dotIndex)
      : fileNameWithExtension;
  const extension =
    dotIndex >= 0 ? fileNameWithExtension.substring(dotIndex + 1) : "";
  return { path, fileName, extension };
};
