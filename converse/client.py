import logging

import grpc

import rpc.converse_pb2 as converse_pb2
import rpc.converse_pb2_grpc as converse_pb2_grpc

logger = logging.getLogger("client")


def run():
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = converse_pb2_grpc.QueryServiceStub(channel)
        response = stub.Query(
            converse_pb2.QueryRequest(query="Am I good leader, how can I do better?")
        )
    logger.info(response.response)


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    run()
