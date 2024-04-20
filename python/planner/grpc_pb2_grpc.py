# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import grpc_pb2 as grpc__pb2


class PlannerStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.CalculateOptimalReplicas = channel.unary_unary(
            '/planner.Planner/CalculateOptimalReplicas',
            request_serializer=grpc__pb2.CalculateOptimalReplicasRequest.
            SerializeToString,
            response_deserializer=grpc__pb2.CalculateOptimalReplicasResponse.
            FromString,
        )


class PlannerServicer(object):
    """Missing associated documentation comment in .proto file."""

    def CalculateOptimalReplicas(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_PlannerServicer_to_server(servicer, server):
    rpc_method_handlers = {
        'CalculateOptimalReplicas':
        grpc.unary_unary_rpc_method_handler(
            servicer.CalculateOptimalReplicas,
            request_deserializer=grpc__pb2.CalculateOptimalReplicasRequest.
            FromString,
            response_serializer=grpc__pb2.CalculateOptimalReplicasResponse.
            SerializeToString,
        ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
        'planner.Planner', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler, ))


# This class is part of an EXPERIMENTAL API.
class Planner(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def CalculateOptimalReplicas(request,
                                 target,
                                 options=(),
                                 channel_credentials=None,
                                 call_credentials=None,
                                 insecure=False,
                                 compression=None,
                                 wait_for_ready=None,
                                 timeout=None,
                                 metadata=None):
        return grpc.experimental.unary_unary(
            request, target, '/planner.Planner/CalculateOptimalReplicas',
            grpc__pb2.CalculateOptimalReplicasRequest.SerializeToString,
            grpc__pb2.CalculateOptimalReplicasResponse.FromString, options,
            channel_credentials, insecure, call_credentials, compression,
            wait_for_ready, timeout, metadata)