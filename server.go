package myrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// MessageReceiver is the interface to receive message from message service
type MessageReceiver interface {
	ReceiveMsg() ([]*RPCMessage, error)
}

// MethodName is a key for map of methods in a service
type MethodName string

// MethodHandler is a abstract handler of a method
type MethodHandler func(service interface{}, in interface{}) (interface{}, error)

// MethodDescription contains method name and its handler
type MethodDescription struct {
	Name    MethodName
	Handler MethodHandler
}

// ServiceName is a key for map of services in a RPC server
type ServiceName string

// ServiceDescription contains all methods and their descriptions
type ServiceDescription struct {
	Name    ServiceName
	Methods map[MethodName]MethodDescription
}

// PayloadDecodeFnc decodes data to output struct
type PayloadDecodeFnc func(data []byte, out interface{}) error

// DefaultQuitSigs are signals that server listen by default
var DefaultQuitSigs = []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}

// RPCServer is struct of this RPC server
type RPCServer struct {
	ctx           context.Context
	locker        sync.Mutex
	servicesList  map[ServiceName]ServiceDescription
	msgReceiver   MessageReceiver
	payloadDecode PayloadDecodeFnc
	exitChan      chan os.Signal
}

// NewRPCServer return new RPC server
func NewRPCServer(ctx context.Context, mr MessageReceiver) *RPCServer {
	return &RPCServer{
		ctx:           ctx,
		msgReceiver:   mr,
		servicesList:  make(map[ServiceName]ServiceDescription),
		payloadDecode: json.Unmarshal, // default
		exitChan:      make(chan os.Signal, 1),
	}
}

// ReplacePayloadDecoder replaces payload decode function of rpc server
func (srv *RPCServer) ReplacePayloadDecoder(decFnc PayloadDecodeFnc) {
	srv.locker.Lock()
	srv.payloadDecode = decFnc
	srv.locker.Unlock()
}

// RegisterService adds new service to server by name and description
func (srv *RPCServer) RegisterService(svName ServiceName, svDesc ServiceDescription) error {
	srv.locker.Lock()
	defer srv.locker.Unlock()
	srv.servicesList[svName] = svDesc
	return nil
}

// ListenQuitSigs listen on signals that make server quit when received.
// This should be called before Serve method
func (srv *RPCServer) ListenQuitSigs(sigs ...os.Signal) {
	if len(sigs) == 0 {
		return
	}

	signal.Notify(srv.exitChan, sigs...)
}

// Serve processes all incoming messages
func (srv *RPCServer) Serve() error {
	defer srv.shutdown()

	for {
		msgs, err := srv.msgReceiver.ReceiveMsg()
		if err != nil {
			return errors.Wrap(err, "error on receive message")
		}

		eg, _ := errgroup.WithContext(srv.ctx)
		for _, msg := range msgs {
			eg.Go(func() error {
				return srv.handleMsg(msg)
			})
		}
		if err := eg.Wait(); err != nil {
			return errors.Wrapf(err, "error on handle message: %+v", msgs)
		}

		select {
		case <-srv.ctx.Done():
			return nil
		case sig := <-srv.exitChan:
			fmt.Println("stop receiving message, got signal: ", sig.String())
			return nil
		default:
		}
	}
}

func (srv *RPCServer) handleMsg(msg *RPCMessage) error {
	svd, ok := srv.servicesList[msg.SvrName]
	if !ok {
		return fmt.Errorf("no service description for %s", msg.SvrName)
	}
	mthd, ok := svd.Methods[msg.MthName]
	if !ok {
		return fmt.Errorf("no method description for %s", msg.MthName)
	}

	payload := srv.payloadDecode(msg.Payload)
	service interface{}, in interface{}
	mthd.Handler()
	return nil
}

func (srv *RPCServer) deleteMsg(msg *RPCMessage) error {
	fmt.Println("deleteMsg ------ ", msg)
	return nil
}

func (srv *RPCServer) shutdown() {
	close(srv.exitChan)
}
