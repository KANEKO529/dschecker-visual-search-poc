import type { RegisterItemEmbeddingResponse } from "../types/embeddingRegistration";

const GENERIC_ERROR_MESSAGE = "Embeddingの登録に失敗しました";

export async function registerItemEmbedding(
  file: File,
  modelNumber: string,
): Promise<RegisterItemEmbeddingResponse> {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;

  const formData = new FormData();
  formData.append("image", file);

  let response: Response;
  try {
    response = await fetch(
      `${baseUrl}/api/poc/items/${encodeURIComponent(modelNumber)}/embeddings`,
      {
        method: "POST",
        body: formData,
      },
    );
  } catch {
    throw new Error(GENERIC_ERROR_MESSAGE);
  }

  if (!response.ok) {
    const body = await response.json().catch(() => null);
    const message =
      body && typeof body.error === "string" ? body.error : GENERIC_ERROR_MESSAGE;
    throw new Error(message);
  }

  return response.json();
}
