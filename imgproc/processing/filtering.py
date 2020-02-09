import io

import pilgram.css
import pilgram
from PIL import Image

import rest


class FilteringServiceImpl(rest.FilteringService):
    def apply(self, src_image_str, filter_name) -> bytes:
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
            raise rest.FilterNotRecognized()
        im = Image.open(io.BytesIO(src_image_str))
        res = getFilterFn[filter_name](im)

        output = io.BytesIO()
        res.save(output, format='JPEG')

        return output.getvalue()
