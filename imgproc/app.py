import logging
import operator
import os
import threading
from argparse import ArgumentParser

from flask import Flask

from imgproc import rest

APP_NAME = "imgproc"
APP_AUTHOR = "midnight coders"
APP_VERSION = "unknown"


def main():
    def setup_logger(dbg):
        logLevel = logging.INFO
        logFormat = "%(asctime)s [%(levelname)s] %(message)s"

        if dbg:
            logLevel = logging.DEBUG
            logFormat = "%(asctime)s [%(levelname)s] %(filename)s: %(message)s"

        logging.basicConfig(level=logLevel, format=logFormat)

    def show_all_routes(server):
        """Display registered routes"""
        rules = []
        for rule in server.url_map.iter_rules():
            methods = ','.join(sorted(rule.methods))
            rules.append((rule.endpoint, methods, str(rule)))

        logging.debug("listing routes:")
        sort_by_rule = operator.itemgetter(2)
        for endpoint, methods, rule in sorted(rules, key=sort_by_rule):
            route = "   {:50s} {:25s} {}".format(endpoint, methods, rule)
            logging.debug(route)

    # parsing cli flags
    parser = ArgumentParser(description="Image processing web-service")
    parser.add_argument(
        "--serviceurl",
        help="url of this web-service in format \"http://<addr>:<port>/\"",
        default=os.environ.get("SERVICEURL", None)
    )
    parser.add_argument(
        "--dbg",
        help="enable debug mode",
        default=os.environ.get("DEBUG", False)
    )

    args = parser.parse_args()
    if not args.serviceurl or len(args.serviceurl.split(':')) < 2:
        exit(parser.print_usage())

    service_url: str = args.serviceurl.split(':')[-2][2:]
    port: int = int(args.serviceurl.split(':')[-1][:-1])

    setup_logger(args.dbg)

    # flask initialization
    app = Flask(__name__)
    app.register_blueprint(rest.chroma_key_controller)

    show_all_routes(app)

    lock = threading.Lock()
    lock.acquire()
    app.run(host=service_url, port=port, debug=args.dbg, use_reloader=False)
    lock.release()


if __name__ == '__main__':
    main()
