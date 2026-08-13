export type SearchItemsByImageResult = {
  modelNumber: string;
  similarity: number;
};

export type SearchItemsByImageResponse = {
  results: SearchItemsByImageResult[];
};
