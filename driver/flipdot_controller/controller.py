# -*- coding: utf-8 -*-
"""Main module."""
import logging
from collections import namedtuple
from typing import Sequence, Union

import numpy as np
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign
from serial import Serial

from flipdot_controller.power import PinConfig, PowerManager

logger = logging.getLogger(__name__)

SignConfig = namedtuple("SignConfig",
                        ['name', 'address', 'width', 'height', 'flip'])

SignInfo = namedtuple("SignInfo", ['name', 'width', 'height'])


class FlipdotController:
    """Adapter class for forwarding server calls into the HanoverController

    Attributes:
        port (Serial): Serial port over which to control the sign
        power_manager (PowerManager): Manager for the power pins
        sign_controller (HanoverController): Low-level controller for the sign
    """

    def __init__(self, port: Serial, signs: Sequence[SignConfig],
                 pins: PinConfig):
        # Create a controller
        self.port = port
        self.sign_controller = HanoverController(self.port)

        # Create signs
        for sign_config in signs:
            config = sign_config._asdict()
            name = config.pop("name")
            sign = HanoverSign(**config)
            self.sign_controller.add_sign(name, sign)

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

    def get_info(self,
                 sign: str = None) -> Union[Sequence[SignInfo], SignInfo]:
        """Return information about the currently attached signs

        Args:
            sign (str, optional): Name of a specific sign to get information on

        Returns:
            Union[Sequence[SignInfo], SignInfo]: Information on the signs
        """
        logger.debug("get_info(sign=%s) called", sign)
        info = {}
        for name, query_sign in self.sign_controller.signs.items():
            info[name] = SignInfo(name=name,
                                  width=query_sign.width,
                                  height=query_sign.height)
        return list(info.values()) if sign is None else info[sign]

    def draw(self, sign: str, image: np.ndarray):
        """Draw the image on the sign

        Args:
            sign (str): The sign to display the image
            image (np.ndarray): The image to display
        """
        logger.debug("draw(sign=%s) called", sign)
        self.sign_controller.draw_image(image, sign)

    def test(self, is_start: bool):
        """Start the test mode on all signs
        """
        logger.debug("test(is_start=%s) called", is_start)
        if is_start:
            self.sign_controller.start_test_signs()
        else:
            self.sign_controller.stop_test_signs()

    def light(self, status: bool):
        """Turn on/off light

        Args:
            status (bool): True to turn on, False to turn off
        """
        logger.debug("light(status=%s) called", status)
        self.power_manager.light(status)
