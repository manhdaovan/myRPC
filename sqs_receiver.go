package myrpc

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/pkg/errors"
)

type sqsReceiver struct {
	ctx      context.Context
	sqs      *sqs.SQS
	conf     ReceiverConf
	queueURL string
}

// NewSQSReceiver returns a SQS client using for receiving message
func NewSQSReceiver(ctx context.Context, conf ReceiverConf) (MessageReceiver, error) {
	sqsClient, queueURL, err := newSQSClient(ctx, conf.Queue)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot init sqs client for sqsReceiver with conf: %+v", conf.Queue)
	}
	return &sqsReceiver{
		sqs:      sqsClient,
		ctx:      ctx,
		conf:     conf,
		queueURL: queueURL,
	}, nil
}

func (sr *sqsReceiver) ReceiveMsg() ([]*RPCMessage, error) {
	param := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sr.queueURL),
		MaxNumberOfMessages: aws.Int64(sr.conf.NumMsgsPerReceive),
		VisibilityTimeout:   aws.Int64(sr.conf.VisibilityTimeout), // sec
		WaitTimeSeconds:     aws.Int64(sr.conf.WaitTimeSeconds),   // sec
	}

	resp, err := sr.sqs.ReceiveMessageWithContext(sr.ctx, param)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot recv message with params: %+v", param)
	}

	ret := make([]*RPCMessage, len(resp.Messages))

	for k, m := range resp.Messages {
		rpcMsg, err := JSONToRPCMsg(*m.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot convert to rpc msg: %+v", m)
		}
		rpcMsg.msgReceiptHandle = *m.ReceiptHandle
		ret[k] = rpcMsg
	}

	return ret, nil
}
