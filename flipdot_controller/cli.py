# -*- coding: utf-8 -*-
"""Console script for flipdot_controller."""
import sys

import click
from serial import Serial

from flipdot_controller.controller import FlipdotController, SignConfig
from flipdot_controller.power import PinConfig
from flipdot_controller.server import Server

SIGNS = [SignConfig(name="top", address=1, width=84, height=7, flip=True)]

PINS = PinConfig(sign=38, light=40)


@click.command()
@click.option('--port', help="Name of serial port")
def main(port):
    """Console script for flipdot_controller."""
    with Serial(port) as ser, FlipdotController(
            port=ser, signs=SIGNS, pins=PINS) as controller:
        server = Server(controller)
        try:
            server.start()
        except KeyboardInterrupt:
            server.stop()


if __name__ == "__main__":
    sys.exit(main())  # pragma: no cover
