"use client";

import { useState, useEffect } from "react";

import { searchItemsByImage } from "../api/searchItemsByImage";
import type { SearchItemsByImageResponse } from "../types/imageSearch";

export default function ImageSearchForm() {
  const [file, setFile] = useState<File | null>(null);
  const [result, setResult] = useState<SearchItemsByImageResponse | null>(
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
    if (!file) {
      return;
    }

    setError(null);
    setResult(null);

    try {
      const response = await searchItemsByImage(file);
      setResult(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : "検索に失敗しました");
    }
  };

  return (
    <div className="w-full max-w-md rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
      <div className="flex flex-col gap-5">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">画像検索</h2>
          <p className="mt-1 text-sm text-gray-500">
            画像を選択して類似商品を検索します。
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

        <button
          type="button"
          onClick={handleSubmit}
          disabled={!file}
          className="w-full rounded-md bg-black px-4 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-80 disabled:cursor-not-allowed disabled:opacity-40"
        >
          検索
        </button>

        {result && (
          <div className="rounded-md bg-gray-50 p-3 text-sm text-gray-700">
            {result.results.length === 0 ? (
              <p>該当する商品が見つかりませんでした。</p>
            ) : (
              <ul className="flex flex-col gap-2">
                {result.results.map((item) => (
                  <li
                    key={item.modelNumber}
                    className="flex items-center justify-between border-b border-gray-200 pb-2 last:border-b-0 last:pb-0"
                  >
                    <span>{item.modelNumber}</span>
                    <span className="text-gray-500">
                      similarity: {item.similarity}
                    </span>
                  </li>
                ))}
              </ul>
            )}
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
