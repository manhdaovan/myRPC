package shared

import (
	"fmt"
)

type GetUserService interface {
	GetUserInfo(in *GetUserRequest) (*GetUserResponse, error)
}

func GetUserServiceExecutor(sv interface{}, requestMsg interface{}) (interface{}, error) {
	in, ok := requestMsg.(*GetUserRequest)
	if !ok {
		return nil, fmt.Errorf("invalid GetUserRequest: %+v", requestMsg)
	}

	return sv.(GetUserService).GetUserInfo(in)
}

func RegisterGetUserService(srv RPCServer, sv GetUserService) {
	var serviceName ServiceName = "GetUserService"
	var serviceDescription ServiceDescription = map[MethodName]MethodExecutor{
		"GetUser": GetUserServiceExecutor,
	}
	srv.AddService(serviceName, serviceDescription)
}
