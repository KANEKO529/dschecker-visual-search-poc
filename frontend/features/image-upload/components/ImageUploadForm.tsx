"use client";

import { useState } from "react";

import { uploadImage } from "../api/uploadImage";
import type { UploadImageResponse } from "../types/imageUpload";

export default function ImageUploadForm() {
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<UploadImageResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

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

      <button
        type="button"
        onClick={handleSubmit}
        disabled={!file}
        className="rounded bg-black px-6 py-2 text-sm text-white transition-opacity hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-40"
      >
        送信
      </button>

      {result && (
        <ul className="text-sm text-gray-600">
          <li>filename: {result.filename}</li>
          <li>size: {result.size}</li>
          <li>contentType: {result.contentType}</li>
        </ul>
      )}

      {error && <p className="text-sm text-red-600">{error}</p>}
    </div>
  );
}