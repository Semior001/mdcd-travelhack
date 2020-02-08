from imgproc.cmd.cli_commander import Commander
from imgproc.cmd.serve import Serve

commands = [
    Serve,
]


def get_commands() -> list[Commander]:
    return commands


__all__ = [
    'get_commands'
]
