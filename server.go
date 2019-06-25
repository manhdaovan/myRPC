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

// MethodName is a key for map of methods in a service
type MethodName string

// MethodHandler is a abstract handler of a method
type MethodHandler func(ctx context.Context, service interface{}, in interface{}) (interface{}, error)

// MethodDecodeFnc decodes data to output struct for current method
type MethodDecodeFnc func(decodeFnc PayloadDecodeFnc, data []byte) (interface{}, error)

// MethodDescription contains method name and its handler
type MethodDescription struct {
	Name          MethodName
	Handler       MethodHandler
	PayloadDecode PayloadDecodeFnc
	DecodeHandle  MethodDecodeFnc
}

// ServiceName is a key for map of services in a RPC server
type ServiceName string

// ServiceDescription contains all methods and their descriptions
type ServiceDescription struct {
	Name    ServiceName
	Methods map[MethodName]MethodDescription
}

// PayloadDecodeFnc decodes bytes into output struct
type PayloadDecodeFnc func(data []byte, out interface{}) error

// DefaultQuitSigs are signals that server listen by default
var DefaultQuitSigs = []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}

// MessageReceiver is the interface to receive message from message service
type MessageReceiver interface {
	ReceiveMsg() ([]*RPCMessage, error)
}

// MessageDeleter is the interface to delete message from message service
type MessageDeleter interface {
	DeleteMsg(msg *RPCMessage) error
}

// RPCServer is struct of this RPC server
type RPCServer struct {
	ctx           context.Context
	locker        sync.Mutex
	servicesDesc  map[ServiceName]ServiceDescription
	services      map[ServiceName]interface{}
	msgReceiver   MessageReceiver
	msgDeleter    MessageDeleter
	payloadDecode PayloadDecodeFnc
	exitChan      chan os.Signal
}

// NewRPCServer return new RPC server
func NewRPCServer(ctx context.Context, mr MessageReceiver, md MessageDeleter) *RPCServer {
	return &RPCServer{
		ctx:           ctx,
		msgReceiver:   mr,
		msgDeleter:    md,
		servicesDesc:  make(map[ServiceName]ServiceDescription),
		services:      make(map[ServiceName]interface{}),
		payloadDecode: json.Unmarshal, // default
		exitChan:      make(chan os.Signal, 1),
	}
}

// ReplacePayloadDecoder replaces payload decode function of rpc server
func (srv *RPCServer) ReplaceDecoder(decFnc PayloadDecodeFnc) {
	srv.locker.Lock()
	srv.payloadDecode = decFnc
	srv.locker.Unlock()
}

// RegisterService adds new service to server by name and description.
// This method SHOULD NOT be called directly outside of RegisterXService() method of each service
func (srv *RPCServer) RegisterService(svc interface{}, svName ServiceName, svDesc ServiceDescription) {
	srv.locker.Lock()
	srv.servicesDesc[svName] = svDesc
	srv.services[svName] = svc
	srv.locker.Unlock()
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
	fmt.Printf("==> handle msg: %p\n", msg)
	// get registered service description from server
	svd, ok := srv.servicesDesc[msg.SvrName]
	if !ok {
		return fmt.Errorf("no service description for %s", msg.SvrName)
	}

	// get registered method description from service
	mthd, ok := svd.Methods[msg.MthName]
	if !ok {
		return fmt.Errorf("no method description for %s", msg.MthName)
	}

	// get registered service instance from server
	svc, ok := srv.services[msg.SvrName]
	if !ok {
		return fmt.Errorf("no service instance for %s", msg.SvrName)
	}

	// decode payload
	decodeFnc := mthd.PayloadDecode
	if decodeFnc == nil {
		// fallback to server default decode func
		decodeFnc = srv.payloadDecode
	}

	in, err := mthd.DecodeHandle(decodeFnc, msg.Payload)
	if err != nil {
		return errors.Wrapf(err, "cannot decode payload msg", msg.Payload)
	}

	_, err = mthd.Handler(srv.ctx, svc, in)
	if err == nil && srv.msgDeleter != nil {
		if err := srv.msgDeleter.DeleteMsg(msg); err != nil {
			return errors.Wrapf(err, "cannot delete msg: %+v", msg)
		}
		fmt.Printf("--> delete msg: %p\n", msg)
	}

	return err
}

func (srv *RPCServer) shutdown() {
	close(srv.exitChan)
}
