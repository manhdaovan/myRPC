package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/manhdaovan/myrpc"
	"github.com/manhdaovan/myrpc/example/message"
	"github.com/manhdaovan/myrpc/example/service"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
)

func main() {
	conf, err := myrpc.SenderConfFromYamlFile("../../config/sender.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on read sender conf file: %+v", err)
		return
	}

	ctx := context.Background()
	sqsSender, err := myrpc.NewSQSSender(ctx, *conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error on init sqs: %+v", err)
		return
	}

	client := myrpc.NewRPCClient(ctx, sqsSender)
	var wg sync.WaitGroup

	// send message in json format
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			msgContent := fmt.Sprintf("Msg to FreeService: %d", idx)
			inMsg := message.FreeMessageIn{Msg: msgContent}
			fmt.Println("send msg: ", msgContent)

			err := client.SendAsyncMsg(service.FreeServiceName, service.FreeServiceEchoMethodName, &inMsg, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error on sending msg to FreeService. msg: %+v, err; %+v", inMsg, err)
			}
		}(i)
	}

	// send message in proto format
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			msgContent := fmt.Sprintf("Msg to ProtoService: %d", idx)
			inMsg := message.EchoProtoIn{Msg: msgContent}
			encFnc := encoding.GetCodec(proto.Name).Marshal
			fmt.Println("send msg: ", msgContent)

			err := client.SendAsyncMsg(service.ProtoServiceName, service.ProtoServiceEchoMethodName, &inMsg, encFnc)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error on sending msg to ProtoService. msg: %+v, err; %+v", inMsg, err)
			}
		}(i)
	}

	wg.Wait()
}
