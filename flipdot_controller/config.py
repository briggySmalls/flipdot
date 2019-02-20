from pathlib import Path

import toml


def _assert(condition, message=None):
    if not condition:
        raise RuntimeError(message)


class ConfigParser:
    def __init__(self, file: Path):
        self.file = file
        self.config = None

    def parse(self):
        self.config = toml.load(self.file)

    def validate(self):
        self._assert_not_missing('serial_port')
        self._assert_not_missing('grpc_port')
        self._assert_not_missing('pins')
        _assert('sign' in self.config['pins'], "sign_pin not supplied in pins")
        _assert('light' in self.config['pins'], "light_pin not supplied in pins")
        self._assert_not_missing('signs')
        _assert(len(self.config['signs']) > 0, "No signs supplied")
        for name, sign in self.config['signs']:
            _assert('address' in sign, "Address missing from sign {}".format(name))
            _assert('width' in sign, "Width missing from sign {}".format(name))
            _assert('height' in sign, "Height missing from sign {}".format(name))

    def config(self):
        self.parse()
        self.validate()
        return self.config

    def _assert_not_missing(self, field):
        _assert(field in self.config, "Config missing: {}".format(field))
