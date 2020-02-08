from imgproc.cmd.cli_commander import Commander
from imgproc.rest import RestParams, ProcessingImageService, Rest

class Serve(Commander):
    def command_name(self) -> str:
        return "serve"

    def execute(self, args):
        Rest(RestParams()).run()