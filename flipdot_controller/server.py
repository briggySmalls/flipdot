# -*- coding: utf-8 -*-
"""Main module."""
from concurrent import futures

import grpc
import numpy as np

from flipdot_controller.controller import FlipdotController
from flipdot_controller.protos.flipdot_pb2 import (DrawResponse, Error,
                                                   GetInfoResponse,
                                                   LightRequest, LightResponse,
                                                   TestRequest, TestResponse, Error)
from flipdot_controller.protos.flipdot_pb2_grpc import (FlipdotServicer,
                                                        add_FlipdotServicer_to_server)


class Server:
    def __init__(self,
                 controller: FlipdotController,
                 max_workers=10,
                 port=5001):
        # Create a servicer
        self.servicer = Servicer(controller)
        # Create gRPC server
        self.server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=max_workers))
        add_FlipdotServicer_to_server(self.servicer, self.server)
        self.server.add_insecure_port('[::]:{}'.format(port))

    def start(self):
        self.server.start()

    def stop(self, grace=0):
        self.server.stop(grace)


class Servicer(FlipdotServicer):
    def __init__(self, controller: FlipdotController):
        self.controller = controller

    def GetInfo(self, request, context) -> GetInfoResponse:
        # Get the sign info
        info = self.controller.get_info()
        # Build a response
        response = GetInfoResponse(error=self._no_error())
        for sign_info in info:
            sign = response.signs.add()
            sign.name = sign_info.name
            sign.width = sign_info.width
            sign.height = sign_info.height
        return response

    def Draw(self, request, context) -> DrawResponse:
        # Determine sign's shape
        sign_info = self.controller.get_info(request.sign)
        # Reconstruct image
        image = np.frombuffer(request.image).reshape((sign_info.height,
                                                      sign_info.width))
        # Send the command
        self.controller.draw(request.sign, image)
        return DrawResponse(error=self._no_error())

    def Test(self, request, context) -> TestResponse:
        if (request.action != TestRequest.START
                and request.action != TestRequest.STOP):
            err = Error(
                code=1, message="Unexpected action {}".format(request.action))
            return TestResponse(error=err)

        self.controller.test(request.action == TestRequest.START)
        return TestResponse(error=self._no_error())

    def Light(self, request, context):
        self.controller.light(request.status == LightRequest.ON)
        return LightResponse(error=self._no_error())

    @staticmethod
    def _no_error():
        return Error(code=0)
