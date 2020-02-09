from flask import request, Blueprint, make_response


class PrintingService:
    def print(self, img_bytes):
        """prints an image to the printer that is connected to the computer"""
        raise NotImplementedError()


class PrintingController:
    service: PrintingService

    def __init__(self, service):
        self.blueprint = Blueprint('chromakeying', __name__)
        self.service = service
        self.routes()

    def routes(self):
        @self.blueprint.route('/print', methods=['POST'])
        def replace_background():
            if 'image' not in request.files:
                return {"error": "image is not provided"}, 400
            image_bytes = request.files['image'].read()
            self.service.print(image_bytes)
            return {
                "ok": True
            }
