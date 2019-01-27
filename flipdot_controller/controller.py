# -*- coding: utf-8 -*-

"""Main module."""
from typing import Sequence
from dataclasses import dataclass, asdict

from serial import Serial
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign
import numpy as np

from flipdot_controller.power import PowerManager, PinConfig


@dataclass
class SignConfig:
    name: str
    address: int
    width: int
    height: int
    flip: bool


@dataclass
class SignInfo:
    name: str
    width: int
    height: int


class FlipdotController:
    def __init__(self, port: Serial, signs: Sequence[SignConfig], power: PinConfig):
        # Create a controller
        self.port = port
        self.sign_controller = HanoverController(self.port)

        # Create signs
        for sign_config in signs:
            sign = HanoverSign(**asdict(sign_config))
            self.sign_controller.add_sign(sign)

        # Create a power manager
        self.power_manager = PowerManager(power)

    def get_info(self) -> Sequence[SignConfig]:
        info = []
        for sign in self.sign_controller._signs.values():
            info.append(SignInfo(name=sign.name, width=sign.width, height=sign.height))
        return info

    def draw(self, sign: str, image: np.ndarray):
        """Draw the image on the sign

        Args:
            sign (str): The sign to display the image
            image (np.ndarray): The image to display
        """
        self.sign_controller.draw_image(image, sign)

    def start_test(self):
        """Start the test mode on all signs
        """
        self.sign_controller.start_test_signs()

    def stop_test(self):
        """Stop the test mode on all signs
        """
        self.sign_controller.stop_test_signs()

    def light(self, status: bool):
        """Turn on/off light

        Args:
            status (bool): True to turn on, False to turn off
        """
        self.power_manager.light(status)
