import platform

from printsrv.rest import PrintingService


class PrintingServiceImpl(PrintingService):
    def __init__(self, local_storage_path='/tmp'):
        self.local_storage_path = local_storage_path

    def print(self, img_bytes):
        """prints an image to the printer that is connected to the computer"""

        platform_name: str = platform.system().capitalize()

        if platform_name == 'LINUX' or platform_name == 'DARWIN':
            import subprocess
            lpr = subprocess.Popen("/usr/bin/lpr", stdin=subprocess.PIPE)
            lpr.stdin.write(img_bytes)
        else:
            import os
            filepath = os.path.join(self.local_storage_path, "tmp.jpg")
            with open(filepath, 'wb') as f:
                f.write(img_bytes)

            os.startfile(filepath, "print")
