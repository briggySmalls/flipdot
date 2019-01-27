# -*- coding: utf-8 -*-

"""Main module."""
from typing import Sequence

from serial import Serial
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign

from flipdot_controller.controller import FlipdotController, PinConfig, SignConfig
from flipdot_controller.protos.flipdot_pb2_grpc import FlipdotServicer
from flipdot_controller.protos.flipdot_pb2 import StartTestResponse, StopTestResponse, LightResponse, DrawResponse, Status


class FlipdotServer(FlipdotServicer):
    def __init__(self, port: Serial, signs: Sequence[SignConfig], power: PinConfig):
        # Create a controller
        self.controller = FlipdotController(port, signs, power)

    def GetInfo(self, request, context):
        pass

    def Draw(self, request, context):
        self.controller.draw(request.sign, request.image)
        return DrawResponse()

    def StartTest(self, request, context):
        self.controller.start_test()
        return StartTestResponse()

    def StopTest(self, request, context):
        self.controller.stop_test()
        return StartTestResponse()

    def Light(self, request, context):
        self.controller.light(request.status == Status.ON)
        return LightResponse()
