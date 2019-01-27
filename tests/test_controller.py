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
    return mock.Mock()


@pytest.fixture
def pins():
    return PinConfig(sign=1, light=2)


def test_get_info(port, pins):
    # Create config for a sign
    sign_config = SignConfig(name='mysign', address=1, width=10, height=8, flip=True)
    # Create the controller
    controller = FlipdotController(port=port, signs=[sign_config], power=pins)
    # Check we get the expected info back
    info = controller.get_info()
    assert len(info) == 1
    assert info[0].name == sign_config.name
    assert info[0].width == sign_config.width
    assert info[0].height == sign_config.height
