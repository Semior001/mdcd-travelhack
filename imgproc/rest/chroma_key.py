import base64

from flask import Blueprint, request, make_response


class ChromaKeyingService:
    def replace(self, src_image_str, bg_image_str) -> str:
        """
        removes background from the image, which is provided as string,
        <b>not filepath</b>
        :param src_image_str: string which represents an image
        :param bg_image_str: string which represents a background image
        :return: string which represents a processed image
        """
        raise NotImplementedError()


class ChromaKeyController:
    def __init__(self, service):
        self.blueprint = Blueprint('chromakeying', __name__)
        self.service = service
        self.routes()

    def routes(self):
        @self.blueprint.route('/replace-background', methods=['POST'])
        def replace_background():
            if 'image' not in request.files or 'background' not in request.files:
                return {"error": "image or background or both are not provided"}, 400
            imageString = request.files['image'].read()
            backgroundString = request.files['background'].read()

            response = make_response(self.service.replace(imageString, backgroundString))
            response.headers.set('Content-Type', 'image/jpeg')
            response.headers.set('Content-Disposition', 'attachment', filename='r.jpg')

            return response
