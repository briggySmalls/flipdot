# -*- coding: utf-8 -*-
"""Console script for flipdot_controller."""
import sys
import time
import logging

import click
from serial import Serial

from flipdot_controller.controller import FlipdotController
from flipdot_controller.server import Server
from flipdot_controller.config import ConfigParser


@click.command()
@click.option('--config', type=click.Path(exists=True), required=True, help="TOML config file")
def main(config):
    """Console script for flipdot_controller."""
    logging.basicConfig(level=logging.DEBUG)
    # Read the config
    parser = ConfigParser.create(config)
    # Create controller from config
    with Serial(parser.basic_config['serial_port']) as ser, FlipdotController(
            port=ser, signs=parser.signs_config, pins=parser.pin_config) as controller:
        server = Server(controller, port=parser.basic_config['grpc_port'])
        try:
            server.start()
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            server.stop()


if __name__ == "__main__":
    sys.exit(main())  # pragma: no cover
