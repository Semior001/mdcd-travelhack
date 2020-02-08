class CommonOpts:
    """
    CommonOpts provides common options to each command
    """
    app_version = None
    app_name = None
    app_author = None

    def __init__(self, app_version, app_name, app_author):
        self.app_name = app_name
        self.app_author = app_author
        self.app_version = app_version


class Commander:
    """
    Commander sets basic common options for each command
    all commands should inherit this class
    """
    common_opts = None

    def __init__(self, common_opts):
        """
        initialize command
        :type common_opts: CommonOpts
        """
        self.common_opts = common_opts

    def command_name(self) -> str:
        raise NotImplementedError("provide command name!")

    def execute(self, args):
        """
        executes the command
        """
        raise NotImplementedError()
