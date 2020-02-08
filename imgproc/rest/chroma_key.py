import base64

from flask import Blueprint, request
from PIL import Image
import io

chroma_key_controller = Blueprint('chromakeying', __name__)


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


@chroma_key_controller.route('/test', methods=['GET', 'POST'])
def test():
    g = request.files.get('imagefile')
    imageString = base64.b64decode(request.form['img'])
    if request.method == 'POST':
        return 'post method'
    else:
        return 'i guess get'


@chroma_key_controller.route('/test1', methods=['PUT'])
def test1():
    return 'PUT method'
