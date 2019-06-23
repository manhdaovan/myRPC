package main

import (
	"context"
	"fmt"
	"os"

	"github.com/manhdaovan/myrpc"
	"github.com/manhdaovan/myrpc/example/service"
)

func main() {
	conf, err := myrpc.ReceiverConfFromYamlFile("../../config/receiver.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on read receiver conf file: %+v", err)
		return
	}

	ctx := context.Background()
	sqsReceiver, err := myrpc.NewSQSReceiver(ctx, *conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on init sqs: %+v", err)
		return
	}

	svr := myrpc.NewRPCServer(ctx, sqsReceiver)
	svr.RegisterService(service.FreeServiceName, service.FreeServiceDes)
	svr.ListenQuitSigs(myrpc.DefaultQuitSigs...)
	if err := svr.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "error on serving: %+v", err)
	}
}
