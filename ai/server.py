import logging
from concurrent import futures

import grpc

import rpc.converse_pb2 as converse_pb2
import rpc.converse_pb2_grpc as converse_pb2_grpc
from command import QueryCommand

logger = logging.getLogger("server")


class QueryService(converse_pb2_grpc.QueryServiceServicer):
    def Query(self, request, context):
        logger.info(f"received query: {request.query}")
        command = QueryCommand()
        response = command.query(request.query)
        return converse_pb2.QueryResponse(response=response)


def serve():
    port = "50051"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    converse_pb2_grpc.add_QueryServiceServicer_to_server(QueryService(), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    logger.info("gRPC server started listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    serve()
