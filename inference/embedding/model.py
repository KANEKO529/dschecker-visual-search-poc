import torch
from PIL import Image
from transformers import AutoImageProcessor, AutoModel

_MODEL_ID = "facebook/dinov3-vits16-pretrain-lvd1689m"
_DEVICE = torch.device("cpu")

_processor: AutoImageProcessor | None = None
_model: AutoModel | None = None


def _load_model() -> tuple[AutoImageProcessor, AutoModel]:
    global _processor, _model
    if _model is None:
        _processor = AutoImageProcessor.from_pretrained(_MODEL_ID)
        _model = AutoModel.from_pretrained(_MODEL_ID).to(_DEVICE)
        _model.eval()
    return _processor, _model


def generate_embedding(image: Image.Image) -> list[float]:
    processor, model = _load_model()
    inputs = processor(images=image, return_tensors="pt").to(_DEVICE)

    with torch.no_grad():
        outputs = model(**inputs)

    return outputs.pooler_output.squeeze(0).tolist()
