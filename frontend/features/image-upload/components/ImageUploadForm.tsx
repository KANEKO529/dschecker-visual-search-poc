"use client";

import { useState, useEffect} from "react";

import { generateEmbedding } from "../api/generateEmbedding";
import { uploadImage } from "../api/uploadImage";
import type {
  GenerateEmbeddingResponse,
  UploadImageResponse,
} from "../types/imageUpload";

const EMBEDDING_PREVIEW_COUNT = 5;

export default function ImageUploadForm() {
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<UploadImageResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [embeddingResult, setEmbeddingResult] =
    useState<GenerateEmbeddingResponse | null>(null);
  const [embeddingError, setEmbeddingError] = useState<string | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  useEffect(() => {
    if (!file) {
      setPreviewUrl(null);
      return;
    }

    const objectUrl = URL.createObjectURL(file);
    setPreviewUrl(objectUrl);

    return () => {
      URL.revokeObjectURL(objectUrl);
    };
  }, [file]);

  const handleSubmit = async () => {
    if (!file) {
      return;
    }

    setError(null);
    setResult(null);

    try {
      const response = await uploadImage(file);
      setResult(response);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "アップロードに失敗しました",
      );
    }
  };

  const handleGenerateEmbedding = async () => {
    if (!file) {
      return;
    }

    setEmbeddingError(null);
    setEmbeddingResult(null);

    try {
      const response = await generateEmbedding(file);
      setEmbeddingResult(response);
    } catch (err) {
      setEmbeddingError(
        err instanceof Error ? err.message : "Embeddingの生成に失敗しました",
      );
    }
  };

  return (
    <div className="flex flex-col items-center gap-6">
      <div className="flex items-center gap-3">
        <label className="cursor-pointer rounded bg-black px-4 py-2 text-sm text-white hover:opacity-80">
          ファイルを選択

          <input
            type="file"
            accept="image/*"
            onChange={(e) => {
              const selectedFile = e.target.files?.[0] ?? null;
              setFile(selectedFile);
            }}
            className="hidden"
          />
        </label>

        <span className="max-w-64 truncate text-sm text-gray-600">
          {file ? file.name : "選択されていません"}
        </span>
      </div>

      {previewUrl && (
        <div className="overflow-hidden rounded border border-gray-200 p-2">
          <img
            src={previewUrl}
            alt="選択した画像のプレビュー"
            className="max-h-80 max-w-full object-contain"
          />
        </div>
      )}

      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={handleSubmit}
          disabled={!file}
          className="rounded bg-black px-6 py-2 text-sm text-white transition-opacity hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-40"
        >
          送信
        </button>

        <button
          type="button"
          onClick={handleGenerateEmbedding}
          disabled={!file}
          className="rounded bg-black px-6 py-2 text-sm text-white transition-opacity hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-40"
        >
          Embeddingを生成
        </button>
      </div>

      {result && (
        <ul className="text-sm text-gray-600">
          <li>filename: {result.filename}</li>
          <li>size: {result.size}</li>
          <li>contentType: {result.contentType}</li>
        </ul>
      )}

      {error && <p className="text-sm text-red-600">{error}</p>}

      {embeddingResult && (
        <ul className="text-sm text-gray-600">
          <li>dimension: {embeddingResult.dimension}</li>
          <li>
            embedding preview:{" "}
            {embeddingResult.embedding
              .slice(0, EMBEDDING_PREVIEW_COUNT)
              .map((value) => value.toFixed(4))
              .join(", ")}
            , ...
          </li>
        </ul>
      )}

      {embeddingError && <p className="text-sm text-red-600">{embeddingError}</p>}
    </div>
  );
}