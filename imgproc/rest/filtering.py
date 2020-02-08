from flask import Blueprint, request, make_response


class FilterNotRecognized(Exception):
    def __init__(self, msg=''):
        self.msg = msg


class FilteringService:
    def apply(self, src_image_str, filter_name) -> bytes:
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
        self.routes()

    def routes(self):
        @self.blueprint.route('/apply-filter', methods=['POST'])
        def apply_filter():
            imageString = request.files['image'].read()
            filter_name = request.form['filter_name']

            if 'filter_name' not in request.form or 'image' not in request.files:
                return {"error": "image or filter name or both are not provided"}, 400

            try:
                repl = self.service.apply(imageString, filter_name)
            except FilterNotRecognized as e:
                return {"error": "cannot recognize filter"}, 400

            response = make_response(repl)
            response.headers.set('Content-Type', 'image/jpeg')
            response.headers.set('Content-Disposition', 'attachment', filename='r.jpg')

            return response
