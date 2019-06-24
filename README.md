# myRPC
My own asynchronous RPC framework that support message in ANY format.
- Default message format is JSON
- Default message transporter is AWS SQS

# Usage
See `/example` directory source code for more details

# What you can custom
- Message encode/decode function
- Message sender/receiver/deleter 

# TODO
- [ ] Auto generate Service code in case of protobuf message
- [ ] Unit test :D

# Future work
- [ ] Support synchronous RPC
