package myrpc

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type sqsSender struct {
	ctx      context.Context
	sqs      *sqs.SQS
	conf     SenderConf
	queueURL string
}

// NewSQSSender returns a SQS client using for sending messages
func NewSQSSender(ctx context.Context, conf SenderConf) (MessageSender, error) {
	sqsClient, queueURL, err := newSQSClient(ctx, conf.Queue)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot init sqs client for sqsSender with conf: %+v", conf.Queue)
	}
	return &sqsSender{
		sqs:      sqsClient,
		ctx:      ctx,
		conf:     conf,
		queueURL: queueURL,
	}, nil
}

// SendAsyncMsg sends message to SQS asynchronously
func (ss *sqsSender) SendAsyncMsg(msg *RPCMessage) error {
	if msg == nil {
		return errors.New("nil msg is given to SendAsyncMsg")
	}

	msgJSON, err := msg.ToJSON()
	if err != nil {
		return errors.Wrapf(err, "cannot convert msg to json: %+v", msg)
	}

	sqsMsg := &sqs.SendMessageInput{
		MessageBody: aws.String(msgJSON),
		QueueUrl:    aws.String(ss.queueURL),
	}
	if _, err := ss.sqs.SendMessageWithContext(ss.ctx, sqsMsg); err != nil {
		return errors.Wrapf(err, "cannot send message to queue: %+v", sqsMsg)
	}

	return nil
}

// SendSyncMsg sends message to SQS, and wait to response.
func (ss *sqsSender) SendSyncMsg(in *RPCMessage, out interface{}) error {
	// TODO: implement waiting response with result payload from SQS
	return errors.New("SQS not support request/reply pattern")
}
