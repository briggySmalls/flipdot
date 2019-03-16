#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Tests for `flipdot_controller` package."""
from unittest import mock

import pytest

# Skip these imports (we have to mock RPi.GPIO first)
from flipdot_controller.controller import FlipdotController, SignConfig  # isort:skip
from flipdot_controller.power import PinConfig  # isort:skip


@pytest.fixture
def port():
    return mock.MagicMock()


@pytest.fixture
def pins():
    return PinConfig(sign=1, light=2)


@pytest.fixture
def controller(pins, port):
    # Create config for a sign
    signs = [
        SignConfig(name='a', address=1, width=10, height=8, flip=True),
        SignConfig(name='b', address=2, width=10, height=8, flip=False)
    ]
    # Create the controller
    return FlipdotController(port=port, signs=signs, pins=pins)


def test_get_info(controller):
    # Check we get the expected info back
    info = controller.get_info()
    assert len(info) == 2
    assert info[0].name == 'a'
    assert info[0].width == 10
    assert info[0].height == 8

    assert info[1].name == 'b'
    assert info[1].width == 10
    assert info[1].height == 8


def test_get_specific_info(controller):
    sign = controller.get_info(sign='b')
    assert sign.name == 'b'
    assert sign.width == 10
    assert sign.height == 8


def test_start_test(controller, port):
    # Send the start command
    controller.test(True)
    # Assert that something was written over serial
    port.write.assert_called_once()


def test_stop_test(controller, port):
    # Send the start command
    controller.test(False)
    # Assert that something was written over serial
    port.write.assert_called_once()


def test_light(controller, pins, mock_rpi):
    # Turn on the lights
    controller.light(True)
    mock_rpi.GPIO.output.assert_called_once_with(pins.light, True)
    mock_rpi.GPIO.reset_mock()
    # Turn off the lights
    controller.light(False)
    mock_rpi.GPIO.output.assert_called_once_with(pins.light, False)
