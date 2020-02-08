#!/usr/bin/env python3
import logging
import os
import sys

from argparse import ArgumentParser

from imgproc.cmd import get_commands

APP_NAME = "imgproc"
APP_AUTHOR = "midnight coders"
APP_VERSION = "unknown"


def main():
    print("{} by {}, version: {}".format(APP_NAME, APP_AUTHOR, APP_VERSION))

    # parsing cli flags
    parser = ArgumentParser(description="Image processing web-service")
    parser.add_argument(
        "command",
        help="command to execute"
    )
    parser.add_argument(
        "--serviceurl",
        help="url of this web-service",
        default=os.environ.get("SERVICEURL", None)
    )
    parser.add_argument(
        "--dbg",
        help="enable debug mode",
        default=os.environ.get("DEBUG", False)
    )

    args = parser.parse_args()

    if not args.serviceurl:
        exit(parser.print_usage())

    setup_logger(args.dbg)

    for cmd in get_commands():
        if args.command == cmd.command_name:
            cmd.execute(sys.argv)


# setups logging level according to the dbg argument
def setup_logger(dbg):
    logLevel = logging.INFO
    logFormat = "%(asctime)s [%(levelname)s] %(message)s"

    if dbg:
        logLevel = logging.DEBUG
        logFormat = "%(asctime)s [%(levelname)s] %(filename)s: %(message)s"

    logging.basicConfig(level=logLevel, format=logFormat)

    logging.debug("Started")
    logging.info("Finished")


if __name__ == "__main__":
    main()
