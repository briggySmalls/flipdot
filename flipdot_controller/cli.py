# -*- coding: utf-8 -*-
"""Console script for flipdot_controller."""
import sys
import time

import click
from serial import Serial

from flipdot_controller.controller import FlipdotController, SignConfig
from flipdot_controller.power import PinConfig
from flipdot_controller.server import Server

SIGNS = [
    SignConfig(name="top", address=1, width=84, height=7, flip=True),
    SignConfig(name="bottom", address=2, width=84, height=7, flip=False),
]

PINS = PinConfig(sign=38, light=40)


@click.command()
@click.option('--serial-port', required=True, help="Name of serial port")
@click.option('--grpc-port', required=True, type=int, help="Number of gRPC port")
def main(serial_port, grpc_port):
    """Console script for flipdot_controller."""
    with Serial(serial_port) as ser, FlipdotController(
            port=ser, signs=SIGNS, pins=PINS) as controller:
        server = Server(controller, port=grpc_port)
        try:
            server.start()
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            server.stop()


if __name__ == "__main__":
    sys.exit(main())  # pragma: no cover
