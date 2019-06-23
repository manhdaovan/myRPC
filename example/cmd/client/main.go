package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/manhdaovan/myrpc"
	"github.com/manhdaovan/myrpc/example/message"
	"github.com/manhdaovan/myrpc/example/service"
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

	var wg sync.WaitGroup
	client := myrpc.NewRPCClient(ctx, sqsSender)
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			msgContent := fmt.Sprintf("Msg %d from client", idx)
			inMsg := message.FreeMessageIn{Msg: msgContent}
			fmt.Println("send msg: ", msgContent)
			err := client.SendAsyncMsg(service.FreeServiceName, service.FreeServiceEchoMethodName, &inMsg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error on SendAsyncMsg: %+v", err)
			}
		}(i)
	}
	wg.Wait()
}
