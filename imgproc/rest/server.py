from flask import request, Flask


class ProcessingImageService:
    """Service that processes images
    """

    def remove_background(self, img_path) -> str:
        """removes background from the given image
        :param img_path: uri to the stored image, that we need to process
        :type img_path: str
        :return: uri to the stored image, that we processed
        """
        raise NotImplementedError()

    def set_background(self, img_path, bg_path) -> str:
        """sets the background to the given image, instead of empty spaces
        :param bg_path: uri to the stored background image
        :type bg_path: str
        :param img_path: uri to the stored image, that we need to process
        :type img_path: str
        :return: uri to the stored image, that we processed
        """
        raise NotImplementedError()


class RestParams:
    app_version = None
    app_name = None
    app_author = None
    service_url = None
    img_proc_service = None

    def __init__(self, app_version, app_name, app_author, service_url,
                 process_service, upload_path):
        """initializes rest params
        :type service_url: str
        :type app_name: str
        :type app_author: str
        :type app_version: str
        :type process_service: ProcessingImageService
        :type upload_path: str
        """
        self.app_version = app_version
        self.app_name = app_name
        self.app_author = app_author
        self.service_url = service_url
        self.img_proc_service = process_service

        # initializing flask
        app = Flask(import_name=__name__)

        app.config['UPLOAD_FOLDER'] = upload_path


class Rest:
    def __init__(self, params):
        """initializing Rest
        :type params: RestParams
        """
        self.params = params

    def run(self):
        """runs the dedicated web-server
        """
        pass
