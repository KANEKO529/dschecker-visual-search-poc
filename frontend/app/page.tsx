import ImageUploadForm from "@/features/image-upload/components/ImageUploadForm";

export default function Home() {
  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-8 font-sans bg-gray-50">
      <main className="flex flex-col items-center gap-2 text-center">
        <h1 className="text-3xl font-semibold tracking-tight text-black">
          DSchecker Visual Search PoC
        </h1>
      </main>
      <ImageUploadForm />
    </div>
  );
}
