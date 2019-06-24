package service

import (
	"context"
	"fmt"

	"github.com/manhdaovan/myrpc"
	"github.com/manhdaovan/myrpc/example/message"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
)

const ProtoServiceName = "ProtoEchoService"
const ProtoServiceEchoMethodName = "ProtoService/Echo"

// ProtoService is a example about processing message in protobuf format
type ProtoService struct{}

// Echo prints incoming message to stdio
func (gs *ProtoService) EchoProto(ctx context.Context, in *message.EchoProtoIn) (*message.EchoProtoOut, error) {
	fmt.Printf("Echo msg from client: %+v\n", in.Msg)
	return nil, nil
}

func RegisterProtoService(svr *myrpc.RPCServer, service EchoProtoServer) {
	svr.RegisterService(service, ProtoServiceName, ProtoServiceDes)
}

// ProtoServiceDes is the description of ProtoService.
var ProtoServiceDes = myrpc.ServiceDescription{
	Name: ProtoServiceName,
	Methods: map[myrpc.MethodName]myrpc.MethodDescription{
		ProtoServiceEchoMethodName: {
			Name:          ProtoServiceEchoMethodName,
			Handler:       ProtoServiceEchoHandler,
			PayloadDecode: encoding.GetCodec(proto.Name).Unmarshal, // use proto Unmarshal as decoder
			DecodeHandle:  ProtoServiceEchoMsgDecode,
		},
	},
}

// ProtoServiceEchoHandler is handler of Echo method
var ProtoServiceEchoHandler = func(ctx context.Context, service, in interface{}) (interface{}, error) {
	svr, ok := service.(EchoProtoServer)
	if !ok {
		return nil, fmt.Errorf("invalid service: %T, %+v", service, service)
	}

	inMsg, ok := in.(*message.EchoProtoIn)
	if !ok {
		return nil, fmt.Errorf("invalid msg: %T, %+v", in, in)
	}

	return svr.EchoProto(ctx, inMsg)
}

// ProtoServiceEchoMsgDecode decodes data to struct to use in ProtoServiceEchoHandler
var ProtoServiceEchoMsgDecode = func(decodeFnc myrpc.PayloadDecodeFnc, data []byte) (interface{}, error) {
	var out message.EchoProtoIn
	err := decodeFnc(data, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
