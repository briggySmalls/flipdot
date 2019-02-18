#!/usr/bin/env python
# -*- coding: utf-8 -*-
"""Tests for `flipdot_controller` package."""
from unittest import mock

import pytest

from flipdot_controller.protos.flipdot_pb2 import LightRequest, TestRequest
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
