from pathlib import Path
from typing import Sequence

import toml

from flipdot_controller.controller import SignConfig
from flipdot_controller.power import PinConfig


def _assert(condition, message=None):
    if not condition:
        raise RuntimeError(message)


class ConfigParser:
    def __init__(self, config: dict):
        self._config = config
        self._validate()

    @staticmethod
    def create(file: Path):
        return ConfigParser(toml.load(str(file)))

    def _validate(self):
        self._assert_not_missing('serial_port')
        self._assert_not_missing('grpc_port')
        self._assert_not_missing('pins')
        _assert('sign' in self._config['pins'],
                "sign_pin not supplied in pins")
        _assert('light' in self._config['pins'],
                "light_pin not supplied in pins")
        self._assert_not_missing('signs')
        _assert(len(self._config['signs']) > 0, "No signs supplied")
        for name, sign in self._config['signs'].items():
            _assert('address' in sign,
                    "Address missing from sign {}".format(name))
            _assert('width' in sign, "Width missing from sign {}".format(name))
            _assert('height' in sign,
                    "Height missing from sign {}".format(name))

    @property
    def basic_config(self):
        return {
            key: value
            for key, value in self._config.items()
            if key not in ['pins', 'signs']
        }

    @property
    def pin_config(self) -> PinConfig:
        return PinConfig(**self._config['pins'])

    @property
    def signs_config(self) -> Sequence[SignConfig]:
        return [
            SignConfig(name=name, **sign_config)
            for name, sign_config in self._config['signs'].items()
        ]

    def _assert_not_missing(self, field):
        _assert(field in self._config, "Config missing: {}".format(field))
