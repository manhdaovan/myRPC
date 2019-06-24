package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/manhdaovan/myrpc"
	"github.com/manhdaovan/myrpc/example/message"
)

// FreeServiceName is used when init RPCMessage
const FreeServiceName = "FreeService"

// FreeServiceEchoMethodName is used when init RPCMessage
const FreeServiceEchoMethodName = "FreeService/Echo"

// FreeServiceI is interface for a non-gRPC service
type FreeServiceI interface {
	Echo(ctx context.Context, in *message.FreeMessageIn) (*message.FreeMessageOut, error)
}

// FreeService is a example about processing message in non-protobuf format
type FreeService struct{}

// Echo prints incoming message to stdio
func (fs *FreeService) Echo(ctx context.Context, in *message.FreeMessageIn) (*message.FreeMessageOut, error) {
	fmt.Printf("Echo msg from client: %+v\n", in.Msg)
	return nil, nil
}

func RegisterFreeService(svr *myrpc.RPCServer, service FreeServiceI) {
	svr.RegisterService(service, FreeServiceName, FreeServiceDes)
}

// FreeServiceDes is the description of FreeService.
var FreeServiceDes = myrpc.ServiceDescription{
	Name: FreeServiceName,
	Methods: map[myrpc.MethodName]myrpc.MethodDescription{
		FreeServiceEchoMethodName: {
			Name:          FreeServiceEchoMethodName,
			Handler:       FreeServiceEchoHandler,
			PayloadDecode: json.Unmarshal, // use json as encoder/decoder
			DecodeHandle:  FreeServiceEchoMsgDecode,
		},
	},
}

// FreeServiceEchoHandler is handler of Echo method
var FreeServiceEchoHandler = func(ctx context.Context, service, in interface{}) (interface{}, error) {
	svr, ok := service.(FreeServiceI)
	if !ok {
		return nil, fmt.Errorf("invalid service: %T, %+v", service, service)
	}

	inMsg, ok := in.(*message.FreeMessageIn)
	if !ok {
		return nil, fmt.Errorf("invalid msg: %T, %+v", in, in)
	}

	return svr.Echo(ctx, inMsg)
}

// FreeServiceEchoMsgDecode decodes data to struct to use in FreeServiceEchoHandler
var FreeServiceEchoMsgDecode = func(decodeFnc myrpc.PayloadDecodeFnc, data []byte) (interface{}, error) {
	var out message.FreeMessageIn
	err := decodeFnc(data, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
