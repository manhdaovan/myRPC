GO_PATH=`go env GOPATH`
PROTO_DIR = ${GO_PATH}/src/github.com/manhdaovan/myrpc/example/proto

install:
	GO111MODULE=on go mod vendor;
	GO111MODULE=on go mod tidy;
gen-proto:
	cd ${PROTO_DIR}; \
	protoc --go_out=plugins=grpc:${GO_PATH}/src ./message/*.proto; \
	protoc --go_out=plugins=grpc:${GO_PATH}/src ./service/*.proto