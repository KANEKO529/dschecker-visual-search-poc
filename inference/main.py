import io

from fastapi import FastAPI, File, UploadFile
from fastapi.responses import JSONResponse
from PIL import Image

from embedding.model import generate_embedding

app = FastAPI()


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/internal/images")
async def receive_image(image: UploadFile | None = File(None)):  # noqa: B008
    if image is None:
        return JSONResponse(status_code=400, content={"error": "image is required"})

    content = await image.read()

    return {
        "filename": image.filename,
        "size": len(content),
        "contentType": image.content_type,
    }


@app.post("/internal/embeddings")
async def create_embedding(image: UploadFile | None = File(None)):  # noqa: B008
    if image is None:
        return JSONResponse(status_code=400, content={"error": "image is required"})

    content = await image.read()
    pil_image = Image.open(io.BytesIO(content)).convert("RGB")
    embedding = generate_embedding(pil_image)

    return {
        "dimension": len(embedding),
        "embedding": embedding,
    }
