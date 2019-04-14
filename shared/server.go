package shared

import "sync"

type ServiceName string
type MethodName string
type MethodExecutor func(service interface{}, in interface{}) (interface{}, error)
type ServiceDescription map[MethodName]MethodExecutor

type RPCServer struct {
	locker       sync.Mutex
	servicesList map[ServiceName]ServiceDescription
}

func (srv *RPCServer) AddService(svName ServiceName, svDesc ServiceDescription) {
	srv.locker.Lock()
	defer srv.locker.Unlock()
	srv.servicesList[svName] = svDesc
}
