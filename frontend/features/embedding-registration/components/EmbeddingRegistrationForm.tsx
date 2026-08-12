"use client";

import { useState, useEffect } from "react";

import { registerItemEmbedding } from "../api/registerItemEmbedding";
import type { RegisterItemEmbeddingResponse } from "../types/embeddingRegistration";

export default function EmbeddingRegistrationForm() {
  const [file, setFile] = useState<File | null>(null);
  const [modelNumber, setModelNumber] = useState("");
  const [result, setResult] = useState<RegisterItemEmbeddingResponse | null>(
    null,
  );
  const [error, setError] = useState<string | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  useEffect(() => {
    return () => {
      if (previewUrl) {
        URL.revokeObjectURL(previewUrl);
      }
    };
  }, [previewUrl]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0] ?? null;
    setFile(selectedFile);
    setPreviewUrl((prev) => {
      if (prev) {
        URL.revokeObjectURL(prev);
      }
      return selectedFile ? URL.createObjectURL(selectedFile) : null;
    });
  };

  const handleSubmit = async () => {
    if (!file || !modelNumber.trim()) {
      return;
    }

    setError(null);
    setResult(null);

    try {
      const response = await registerItemEmbedding(file, modelNumber);
      setResult(response);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Embeddingの登録に失敗しました",
      );
    }
  };

  return (
    <div className="w-full max-w-md rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="flex flex-col gap-5">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">
            Embedding登録
          </h2>
          <p className="mt-1 text-sm text-gray-500">
            画像と型番を指定してEmbeddingを登録します。
          </p>
        </div>

        <div>
          <label className="mb-2 block text-sm font-medium text-gray-700">
            画像
          </label>

          <div className="flex items-center gap-3">
            <label className="cursor-pointer rounded-md bg-black px-4 py-2 text-sm text-white hover:opacity-80">
              ファイルを選択

              <input
                type="file"
                accept="image/*"
                onChange={handleFileChange}
                className="hidden"
              />
            </label>

            <span className="min-w-0 flex-1 truncate text-sm text-gray-500">
              {file ? file.name : "選択されていません"}
            </span>
          </div>
        </div>

        {previewUrl && (
          <div className="flex justify-center rounded-lg border border-gray-200 bg-gray-50 p-3">
            <img
              src={previewUrl}
              alt="選択した画像のプレビュー"
              className="max-h-64 max-w-full object-contain"
            />
          </div>
        )}

        <div>
          <label
            htmlFor="modelNumber"
            className="mb-2 block text-sm font-medium text-gray-700"
          >
            型番
          </label>

          <input
            id="modelNumber"
            type="text"
            value={modelNumber}
            onChange={(e) => setModelNumber(e.target.value)}
            placeholder="例: NTR-A5RJ-JPN"
            className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-black outline-none focus:border-black"
          />
        </div>

        <button
          type="button"
          onClick={handleSubmit}
          disabled={!file || !modelNumber.trim()}
          className="w-full rounded-md bg-black px-4 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-40"
        >
          Embeddingを登録
        </button>

        {result && (
          <div className="rounded-md bg-gray-50 p-3 text-sm text-gray-700">
            <p>登録に成功しました。</p>
            <p className="mt-2">型番: {result.modelNumber}</p>
            <p>Embedding ID: {result.embeddingId}</p>
          </div>
        )}

        {error && (
          <div className="rounded-md bg-red-50 p-3 text-sm text-red-600">
            {error}
          </div>
        )}
      </div>
    </div>
  );
}
