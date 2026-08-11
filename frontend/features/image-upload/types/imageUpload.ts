export type UploadImageResponse = {
  filename: string;
  size: number;
  contentType: string;
};

export type GenerateEmbeddingResponse = {
  dimension: number;
  embedding: number[];
};
