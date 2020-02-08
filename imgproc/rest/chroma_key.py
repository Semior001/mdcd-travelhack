from flask import Blueprint, request

chroma_key_controller = Blueprint('chromakeying', __name__)


@chroma_key_controller.route('/test', methods=['GET', 'POST'])
def test():
    if request.method == 'POST':
        return 'post method'
    else:
        return 'i guess get'


@chroma_key_controller.route('/test1', methods=['PUT'])
def test1():
    return 'PUT method'
