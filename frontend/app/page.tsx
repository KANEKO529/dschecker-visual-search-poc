import EmbeddingRegistrationForm from "@/features/embedding-registration/components/EmbeddingRegistrationForm";
import ImageUploadForm from "@/features/image-upload/components/ImageUploadForm";

export default function Home() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center gap-8 bg-gray-50 px-6 py-10 font-sans">
      <main className="text-center">
        <h1 className="text-3xl font-semibold tracking-tight text-black">
          DSchecker Visual Search PoC
        </h1>
      </main>

      <div className="flex w-full max-w-5xl flex-col items-start justify-center gap-6 md:flex-row">
        <ImageUploadForm />
        <EmbeddingRegistrationForm />
      </div>
    </div>
  );
}
