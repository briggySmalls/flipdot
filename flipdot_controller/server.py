# -*- coding: utf-8 -*-
"""Main module."""
from concurrent import futures
from typing import Sequence

import grpc
from pyflipdot.pyflipdot import HanoverController
from pyflipdot.sign import HanoverSign
from serial import Serial

from flipdot_controller.controller import (FlipdotController, PinConfig,
                                           SignConfig)
from flipdot_controller.protos.flipdot_pb2 import (DrawResponse, LightRequest,
                                                   LightResponse,
                                                   StartTestResponse,
                                                   StopTestResponse)
from flipdot_controller.protos.flipdot_pb2_grpc import (FlipdotServicer,
                                                        add_FlipdotServicer_to_server)


class Server:
    def __init__(self,
                 controller: FlipdotController,
                 max_workers=10,
                 port=50051):
        # Create a servicer
        self.servicer = FlipdotServicer(controller)
        # Create gRPC server
        self.server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=max_workers))
        add_FlipdotServicer_to_server(self, self.server)
        self.server.add_insecure_port('[::]:{}'.format(port))

    def start(self):
        self.server.start()

    def stop(self):
        self.server.stop()


class Servicer(FlipdotServicer):
    def __init__(self, controller: FlipdotController):
        self.controller = controller

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
        self.controller.light(request.status == LightRequest.Status.ON)
        return LightResponse()
