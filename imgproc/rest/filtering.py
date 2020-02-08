import base64

from flask import Blueprint, request
from PIL import Image
import pilgram
import pilgram.css


class FilteringService:
    def apply(self, src_image_str, filter_name) -> str:
        """
        apply filter to the image, which is provided as string,
        <b>not filepath</b>
        :param src_image_str: string which represents an image
        :param filter_name: string which represents a filter
        :return: string which represents a processed image
        """
        getFilterFn = {
            "_1977": pilgram._1977,
            "aden": pilgram.aden,
            "brannan": pilgram.brannan,
            "brooklyn": pilgram.brooklyn,
            "clarendon": pilgram.clarendon,
            "earlybird": pilgram.earlybird,
            "gingham": pilgram.gingham,
            "hudson": pilgram.hudson,
            "inkwell": pilgram.inkwell,
            "kelvin": pilgram.kelvin,
            "lark": pilgram.lark,
            "lofi": pilgram.lofi,
            "maven": pilgram.maven,
            "mayfair": pilgram.mayfair,
            "moon": pilgram.moon,
            "nashville": pilgram.nashville,
            "perpetua": pilgram.perpetua,
            "reyes": pilgram.reyes,
            "rise": pilgram.rise,
            "slumber": pilgram.slumber,
            "stinson": pilgram.stinson,
            "toaster": pilgram.toaster,
            "valencia": pilgram.valencia,
            "walden": pilgram.walden,
            "willow": pilgram.willow,
            "xpro2": pilgram.xpro2,
            "css-contrast": pilgram.css.contrast,
            "css-grayscale": pilgram.css.grayscale,
            "css-hue_rotate": pilgram.css.hue_rotate,
            "css-saturate": pilgram.css.saturate,
            "css-sepia": pilgram.css.sepia,
        }
        if filter_name not in getFilterFn.keys():
            raise ValueError  # TODO
        im = Image.open(src_image_str)
        getFilterFn[filter_name](im).save()
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
