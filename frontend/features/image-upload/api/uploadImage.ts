import type { UploadImageResponse } from "../types/imageUpload";

const GENERIC_ERROR_MESSAGE = "アップロードに失敗しました";

export async function uploadImage(file: File): Promise<UploadImageResponse> {
  const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL;

  const formData = new FormData();
  formData.append("image", file);

  let response: Response;
  try {
    response = await fetch(`${baseUrl}/api/poc/images`, {
      method: "POST",
      body: formData,
    });
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
