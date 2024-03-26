python -m grpc_tools.protoc --proto_path=protos \
    --python_out=. \
    --grpc_python_out=. \
    --pyi_out=. \
    protos/converse.proto
