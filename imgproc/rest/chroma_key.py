import base64

from flask import Blueprint, request
from PIL import Image
import io
import cv2
import numpy as np
import matplotlib.pyplot as plt


class ChromaKeyingService:
    def replace(self, src_image_str, bg_image_str) -> str:
        """
        removes background from the image, which is provided as string,
        <b>not filepath</b>
        :param src_image_str: string which represents an image
        :param bg_image_str: string which represents a background image
        :return: string which represents a processed image
        """
        img = cv2.imread(src_image_str)
        RED, GREEN, BLUE = (2, 1, 0)
        reds = img[:, :, RED]
        greens = img[:, :, GREEN]
        blues = img[:, :, BLUE]
        mask = (greens < 70) | (reds > greens) | (blues > greens)
        result = np.where(mask, 255, 0)
        # TODO background replacing
        # cv2.imwrite("./test.jpg", result)
        raise NotImplementedError()


class ChromaKeyController:
    def __init__(self, service):
        self.blueprint = Blueprint('chromakeying', __name__)
        self.service = service

    def routes(self):
        @self.blueprint.route('/replace-background', methods=['GET', 'POST'])
        def replace_background():
            imageString = base64.b64decode(request.form['image'])
            backgroundString = base64.b64decode(request.form['background'])
            if request.method == 'POST':
                return self.service.replace(imageString, backgroundString)
