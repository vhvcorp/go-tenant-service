# Generated Protocol Buffer Files

This directory contains auto-generated code from `.proto` files.

## How to Generate

Run the following command from the project root:

```bash
# For tenant-service
cd server/proto
make proto

# Or use the framework-wide script
cd ../../go-framework/scripts
.\generate-proto.bat
```

## Files

After generation, you will find:
- `tenant.pb.go` - Protocol buffer message definitions
- `tenant_grpc.pb.go` - gRPC service client and server stubs
- `tenant.pb.gw.go` - gRPC-Gateway reverse proxy

**Do not edit these files manually. They are auto-generated.**

## Prerequisites

- protoc (Protocol Buffers compiler)
- protoc-gen-go
- protoc-gen-go-grpc
- protoc-gen-grpc-gateway

See `../../../go-framework/scripts/PROTO_SETUP.md` for installation instructions.
