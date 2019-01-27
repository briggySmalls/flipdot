from dataclasses import dataclass

import RPi.GPIO as GPIO


@dataclass
class PinConfig:
    sign: int
    light: int


class PowerManager(object):
    def __init__(self, pins: PinConfig):
        GPIO.setmode(GPIO.BOARD)
        self.pins = pins

        # Configure the pins as outputs
        GPIO.setup(self.pins.sign, GPIO.OUT, initial=GPIO.HIGH)
        GPIO.setup(self.pins.light, GPIO.OUT, initial=GPIO.HIGH)

    def __enter__(self):
        return self

    def __exit__(self, *args, **kwargs):
        # Cleanup GPIO usage
        GPIO.cleanup()

    def light(self, status: bool):
        PowerManager._write_pin(self.pins.light, status)

    def sign(self, status: bool):
        PowerManager._write_pin(self.pins.sign, status)

    @staticmethod
    def _write_pin(pin: int, status: bool):
        GPIO.output(pin, not status)
