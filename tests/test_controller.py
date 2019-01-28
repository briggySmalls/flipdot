#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""Tests for `flipdot_controller` package."""
from unittest import mock
from dataclasses import asdict

import pytest

mock_rpi = mock.MagicMock()
modules = {
    "RPi": mock_rpi,
    "RPi.GPIO": mock_rpi.GPIO,
}
patcher = mock.patch.dict("sys.modules", modules)
patcher.start()

from flipdot_controller.power import PinConfig
from flipdot_controller.controller import FlipdotController, SignConfig


@pytest.fixture
def port():
    return mock.MagicMock()


@pytest.fixture
def pins():
    return PinConfig(sign=1, light=2)


@pytest.fixture
def controller(pins, port):
    # Create config for a sign
    sign_config = SignConfig(name='mysign', address=1, width=10, height=8, flip=True)
    # Create the controller
    return FlipdotController(port=port, signs=[sign_config], power=pins)


def test_get_info(controller):
    # Check we get the expected info back
    info = controller.get_info()
    assert len(info) == 1
    assert info[0].name == 'mysign'
    assert info[0].width == 10
    assert info[0].height == 8


def test_start_test(controller, port):
    # Send the start command
    controller.start_test()
    # Assert that something was written over serial
    port.write.assert_called_once()

def test_stop_test(controller, port):
    # Send the start command
    controller.stop_test()
    # Assert that something was written over serial
    port.write.assert_called_once()
