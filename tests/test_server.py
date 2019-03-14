#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Tests for `flipdot_controller` package."""
from unittest import mock

import numpy as np
import pytest

from flipdot_controller.controller import SignInfo
from flipdot_controller.protos.flipdot_pb2 import (DrawRequest, Image, LightRequest,
                                                   TestRequest)
from flipdot_controller.server import Servicer


@pytest.fixture
def controller():
    return mock.MagicMock()


@pytest.fixture
def servicer(controller):
    # Create the servicer
    return Servicer(controller)


def test_start_test(servicer, controller):
    # Test starting
    request = TestRequest(action=TestRequest.START)
    servicer.Test(request, None)
    controller.test.assert_called_once_with(True)


def test_stop_test(servicer, controller):
    # Test stopping
    request = TestRequest(action=TestRequest.STOP)
    servicer.Test(request, None)
    controller.test.assert_called_once_with(False)


def test_light_on(servicer, controller):
    # Send the start command
    request = LightRequest(status=LightRequest.ON)
    servicer.Light(request, None)
    controller.light.assert_called_with(True)


def test_light_off(servicer, controller):
    # Send the start command
    request = LightRequest(status=LightRequest.OFF)
    servicer.Light(request, None)
    controller.light.assert_called_with(False)


def test_get_info(servicer, controller):
    # Test getting some signs
    test_input = [
        SignInfo(name='test1', width=3, height=2),
        SignInfo(name='test2', width=3, height=2),
    ]
    controller.get_info.return_value = test_input
    # Send the request
    response = servicer.GetInfo(None, None)
    # Assert the response
    assert len(response.signs) == 2
    for expected, actual in zip(test_input, response.signs):
        assert actual.name == expected.name
        assert actual.width == expected.width
        assert actual.height == expected.height


def test_draw(servicer, controller):
    # Create a fake sign to return
    controller.get_info.return_value = SignInfo(name='test', width=3, height=2)
    # Create a draw request
    img = [False, False, True, False, False, True]
    request = DrawRequest(sign='test', image=Image(data=img))
    # Send the request
    servicer.Draw(request, None)
    controller.get_info.assert_called_once_with('test')
    controller.draw.assert_called_once()
    args, _ = controller.draw.call_args
    called_sign, called_image = args
    assert called_sign == 'test'
    np.testing.assert_equal(called_image,
                            [[False, False, True], [False, False, True]])
