import logging

import grpc

import converse_pb2
import converse_pb2_grpc

logger = logging.getLogger(__name__)


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = converse_pb2_grpc.QueryServiceStub(channel)
        response = stub.Query(converse_pb2.QueryRequest(query="Hello, server!"))
    logger.info("Client received: " + response.response)


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    run()
