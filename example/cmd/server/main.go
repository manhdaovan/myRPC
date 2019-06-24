package main

import (
	"context"
	"fmt"
	"os"

	"github.com/manhdaovan/myrpc"
	"github.com/manhdaovan/myrpc/example/service"
)

func main() {
	ctx := context.Background()

	rconf, err := myrpc.ReceiverConfFromYamlFile("../../config/receiver.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on read receiver conf file: %+v", err)
		return
	}
	sqsReceiver, err := myrpc.NewSQSReceiver(ctx, *rconf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on init sqsReceiver: %+v", err)
		return
	}

	dconf, err := myrpc.DeleterConfFromYamlFile("../../config/deleter.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on read deleter conf file: %+v", err)
		return
	}
	sqsDeleter, err := myrpc.NewSQSDeleter(ctx, *dconf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on init sqsDeleter: %+v", err)
		return
	}

	// init server
	svr := myrpc.NewRPCServer(ctx, sqsReceiver, sqsDeleter)

	// register all services to server
	service.RegisterFreeService(svr, &service.FreeService{})
	service.RegisterProtoService(svr, &service.ProtoService{})

	svr.ListenQuitSigs(myrpc.DefaultQuitSigs...)
	fmt.Println("===== server started ===== ")
	if err := svr.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "error on serving: %+v", err)
	}
}
