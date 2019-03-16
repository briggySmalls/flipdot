"""Controls power GPIOs"""

from collections import namedtuple

import RPi.GPIO as GPIO  # pylint: disable=E0401

PinConfig = namedtuple("PinConfig", ['sign', 'light'])


class PowerManager:
    """Manager for providing power to components

    Attributes:
        pins (PinConfig): GPIO name/pin mapping
    """

    def __init__(self, pins: PinConfig):
        GPIO.setmode(GPIO.BOARD)
        self.pins = pins

        # Configure the pins as outputs
        GPIO.setup(self.pins.sign, GPIO.OUT, initial=GPIO.LOW)
        GPIO.setup(self.pins.light, GPIO.OUT, initial=GPIO.LOW)

    def __enter__(self):
        return self

    def __exit__(self, *args, **kwargs):
        # Cleanup GPIO usage
        GPIO.cleanup()

    def light(self, status: bool):
        """Set power to lights

        Args:
            status (bool): True if on, otherwise False
        """
        PowerManager._write_pin(self.pins.light, status)

    def sign(self, status: bool):
        """Set power to signs

        Args:
            status (bool): True if on, otherwise False
        """
        PowerManager._write_pin(self.pins.sign, status)

    @staticmethod
    def _write_pin(pin: int, status: bool):
        GPIO.output(pin, status)
