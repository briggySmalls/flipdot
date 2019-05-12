# -*- coding: utf-8 -*-
"""Main module."""
import logging
from concurrent import futures

import grpc
import numpy as np
from grpc_reflection.v1alpha import reflection

from flipdot_controller.controller import FlipdotController
from flipdot_controller.protos.driver_pb2 import (DESCRIPTOR, DrawResponse,
                                                   GetInfoResponse,
                                                   LightRequest, LightResponse,
                                                   TestRequest, TestResponse)
from flipdot_controller.protos.driver_pb2_grpc import (
    DriverServicer, add_DriverServicer_to_server)

logger = logging.getLogger(__name__)


class Server:
    """Helper for starting gRPC server with Flipdot service

    Attributes:
        server (grpc.Server): The gRPC server
        servicer (Servicer): The Flipdot servicer
    """

    def __init__(self,
                 controller: FlipdotController,
                 max_workers=10,
                 port=5001):
        # Create a servicer
        self.servicer = Servicer(controller)
        # Create gRPC server
        self.server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=max_workers))
        add_DriverServicer_to_server(self.servicer, self.server)
        # the reflection service will be aware of "Flipdot" and
        # "ServerReflection" services.
        service_names = (
            DESCRIPTOR.services_by_name['Driver'].full_name,
            reflection.SERVICE_NAME,
        )
        reflection.enable_server_reflection(service_names, self.server)
        port_string = '[::]:{}'.format(port)
        logger.debug("Starting server on port '%s'", port_string)
        self.server.add_insecure_port(port_string)

    def start(self):
        """Starts the server listening
        """
        self.server.start()

    def stop(self, grace=0):
        """Stops the server

        Args:
            grace (int, optional): A duration of time in seconds or None.
        """
        self.server.stop(grace)


class Servicer(DriverServicer):
    """Servicers for the Driver service
    Generally just forwards on calls to the controller, and creates appropriate
    error codes

    Attributes:
        controller (FlipdotController): The controller to forward calls to
    """

    def __init__(self, controller: FlipdotController):
        self.controller = controller

    def GetInfo(self, request, context) -> GetInfoResponse:
        # Get the sign info
        info = self.controller.get_info()
        # Build a response
        response = GetInfoResponse()
        for sign_info in info:
            sign = response.signs.add()  # pylint: disable=E1101
            sign.name = sign_info.name
            sign.width = sign_info.width
            sign.height = sign_info.height
        return response

    def Draw(self, request, context) -> DrawResponse:
        # Determine sign's shape
        sign_info = self.controller.get_info(request.sign)
        # Reconstruct image
        image = np.array(request.image.data, dtype=bool).reshape(
            (sign_info.height, sign_info.width))
        # Send the command
        self.controller.draw(request.sign, image)
        return DrawResponse()

    def Test(self, request, context) -> TestResponse:
        if (request.action != TestRequest.START  # noqa pylint: disable=E1101
                and request.action != TestRequest.STOP):  # noqa pylint: disable=E1101
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("Unexpected action {}".format(request.action))
            return TestResponse()

        self.controller.test(request.action == TestRequest.START)  # noqa pylint: disable=E1101
        return TestResponse()

    def Light(self, request, context):
        if (request.status != LightRequest.ON  # noqa pylint: disable=E1101
                and request.status != LightRequest.OFF):  # noqa pylint: disable=E1101
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details("Unexpected status {}".format(request.status))
            return LightResponse()

        self.controller.light(request.status == LightRequest.ON)  # noqa pylint: disable=E1101
        return LightResponse()
