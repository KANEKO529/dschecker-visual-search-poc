from fastapi import FastAPI, File, UploadFile
from fastapi.responses import JSONResponse

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
