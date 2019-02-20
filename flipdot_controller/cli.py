# -*- coding: utf-8 -*-
"""Console script for flipdot_controller."""
import sys
import time
import logging

import click
from serial import Serial

from flipdot_controller.controller import FlipdotController, SignConfig
from flipdot_controller.power import PinConfig
from flipdot_controller.server import Server
from flipdot_controller.config import ConfigParser

SIGNS = [
    SignConfig(name="top", address=1, width=84, height=7, flip=True),
    SignConfig(name="bottom", address=2, width=84, height=7, flip=False),
]


@click.command()
@click.option('--config', type=click.Path(exists=True), required=True, help="TOML config file")
def main(config):
    """Console script for flipdot_controller."""
    logging.basicConfig(level=logging.DEBUG)
    # Read the config
    config = ConfigParser.create(config)
    # Create controller from config
    with Serial(config.config['serial_port']) as ser, FlipdotController(
            port=ser, signs=config.signs_config, pins=config.pin_config) as controller:
        server = Server(controller, port=config['grpc_port'])
        try:
            server.start()
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            server.stop()


if __name__ == "__main__":
    sys.exit(main())  # pragma: no cover
