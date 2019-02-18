# -*- coding: utf-8 -*-
"""Main module."""
from collections import namedtuple
from typing import Sequence, Union

import numpy as np
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign
from serial import Serial

from flipdot_controller.power import PinConfig, PowerManager

SignConfig = namedtuple("SignConfig",
                        ['name', 'address', 'width', 'height', 'flip'])

SignInfo = namedtuple("SignInfo", ['name', 'width', 'height'])


class FlipdotController:
    def __init__(self, port: Serial, signs: Sequence[SignConfig],
                 pins: PinConfig):
        # Create a controller
        self.port = port
        self.sign_controller = HanoverController(self.port)

        # Create signs
        for sign_config in signs:
            sign = HanoverSign(**sign_config._asdict())
            self.sign_controller.add_sign(sign)

        # Create a power manager
        self.power_manager = PowerManager(pins)

    def __enter__(self):
        # Turn on the sign
        self.power_manager.sign(True)
        return self

    def __exit__(self, *args, **kwargs):
        # Turn off the sign
        self.power_manager.sign(False)
        # Cleanup GPIOs
        self.power_manager.__exit__(*args, **kwargs)

    def get_info(self, sign=None) -> Union[Sequence[SignInfo], SignInfo]:
        info = {}
        for s in self.sign_controller._signs.values():
            info[s.name] = SignInfo(
                name=s.name, width=s.width, height=s.height)
        return list(info.values()) if sign is None else info[sign]

    def draw(self, sign: str, image: np.ndarray):
        """Draw the image on the sign

        Args:
            sign (str): The sign to display the image
            image (np.ndarray): The image to display
        """
        self.sign_controller.draw_image(image, sign)

    def test(self, is_start: bool):
        """Start the test mode on all signs
        """
        if is_start:
            self.sign_controller.start_test_signs()
        else:
            self.sign_controller.stop_test_signs()

    def light(self, status: bool):
        """Turn on/off light

        Args:
            status (bool): True to turn on, False to turn off
        """
        self.power_manager.light(status)
