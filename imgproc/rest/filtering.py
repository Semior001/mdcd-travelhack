import base64

from flask import Blueprint, request
from PIL import Image
import io


class FilteringService:
    def apply(self, src_image_str, filter_name) -> str:
        """
        apply filter to the image, which is provided as string,
        <b>not filepath</b>
        :param src_image_str: string which represents an image
        :param filter_name: string which represents a filter
        :return: string which represents a processed image
        """
        raise NotImplementedError()


class FilterController:
    def __init__(self, service):
        self.blueprint = Blueprint('filtering', __name__)
        self.service = service

    def routes(self):
        @self.blueprint.route('/apply-filter', methods=['GET', 'POST'])
        def apply_filter():
            imageString = base64.b64decode(request.form['image'])
            filter_name = request.form['filter_name']
            if request.method == 'POST':
                return self.service.replace(imageString, filter_name)
