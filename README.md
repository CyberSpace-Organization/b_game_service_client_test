To compile the protobuf file, use

```shell
    protoc -I proto/ proto/room.proto --go_out=plugins=grpc:go  
```