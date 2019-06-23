package service

import (
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
	Echo(in *message.FreeMessageIn) (*message.FreeMessageOut, error)
}

// FreeService is a example about non-gRPC service
type FreeService struct{}

// Echo prints incoming message to stdio
func (fs *FreeService) Echo(in *message.FreeMessageIn) (*message.FreeMessageOut, error) {
	fmt.Printf("Echo msg: %+v\n", in)
	return nil, nil
}

// FreeServiceDes is the description of FreeService.
var FreeServiceDes = myrpc.ServiceDescription{
	Name: FreeServiceName,
	Methods: map[myrpc.MethodName]myrpc.MethodDescription{
		FreeServiceEchoMethodName: myrpc.MethodDescription{
			Name:    FreeServiceEchoMethodName,
			Handler: FreeServiceEchoHandler,
		},
	},
}

// FreeServiceEchoHandler is handler of Echo method
var FreeServiceEchoHandler = func(service, in interface{}) (interface{}, error) {
	svr, ok := service.(FreeServiceI)
	if !ok {
		return nil, fmt.Errorf("invalid service: %T, %+v", service, service)
	}

	inMsg, ok := in.(message.FreeMessageIn)
	if !ok {
		return nil, fmt.Errorf("invalid msg: %T, %+v", in, in)
	}

	return svr.Echo(&inMsg)
}
