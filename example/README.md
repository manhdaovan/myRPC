# How to run
- Install protobuf: `brew install protobuf && brew upgrade protobuf`
- Install dependencies: `make install`
- Start elasticmq(a msg service compatible SQS interfaces): `docker-compose up`
- Open new terminal tab and start server: `cd cmd/server && go run main.go`
- Open new terminal tab and start client: `cd cmd/client && go run main.go`