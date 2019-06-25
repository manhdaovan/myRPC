# Write your own RPC framework, why not?
https://speakerdeck.com/manhdaovan/write-your-own-async-rpc-framework

# myRPC
My own asynchronous RPC framework that support message in ANY format.
- Default message format is JSON. Support protobuf message without any touch.
- Default message transporter is AWS SQS.

# Usage
- If you use message in JSON format:
  - Define message struct as in `example/message/free_message.go`
  - Define service that having RPC interfaces and implementation as in `example/service/free_service.go`
    - The service and methods descriptions part could be auto generate after having message struct and RPC interfaces
  - Register it on server side as in `example/cmd/server/main.go`
- If you use message in Protobuf format:
  - Generate message struct and RPC interfaces using `protoc` as in `example/Makefile`
  - Define service and method descriptions as in `example/service/grpc_service_extend.go`, and implement RPC interfaces
    - The service and methods descriptions part could be auto generate after having message struct and RPC interfaces
  - Register it on server side as in `example/cmd/server/main.go`
- Or you can create your own message format and use it, by implement your own message encoder/decoder
  - Other steps are similar to above steps

# Example
See `/example` directory source code for more details

# What you can custom
- Message encode/decode function
- Message sender/receiver/deleter 

# TODO
- [ ] Auto generate Service code in case of protobuf message
- [ ] Unit test :D

# Future work
- [ ] Support synchronous RPC
