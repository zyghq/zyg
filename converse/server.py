import logging
from concurrent import futures

import grpc

import converse_pb2
import converse_pb2_grpc

logger = logging.getLogger(__name__)


class QueryService(converse_pb2_grpc.QueryServiceServicer):
    def Query(self, request, context):
        logger.info(f"Received query: {request.query}")
        logger.info("do some backend magic here...")
        return converse_pb2.QueryResponse(response="Hello, you said: " + request.query)


def serve():
    port = "50051"
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    converse_pb2_grpc.add_QueryServiceServicer_to_server(QueryService(), server)
    server.add_insecure_port("[::]:" + port)
    server.start()
    logger.info("Server started, listening on " + port)
    server.wait_for_termination()


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO)
    serve()
