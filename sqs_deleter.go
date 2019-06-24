package myrpc

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type sqsDeleter struct {
	ctx      context.Context
	sqs      *sqs.SQS
	conf     DeleterConf
	queueURL string
}

// NewSQSDeleter returns a SQS client using for deleting message
func NewSQSDeleter(ctx context.Context, conf DeleterConf) (MessageDeleter, error) {
	sqsClient, queueURL, err := newSQSClient(ctx, conf.Queue)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot init sqs client for sqsDeleter with conf: %+v", conf.Queue)
	}

	return &sqsDeleter{
		sqs:      sqsClient,
		ctx:      ctx,
		conf:     conf,
		queueURL: queueURL,
	}, nil
}

func (sr *sqsDeleter) DeleteMsg(msg *RPCMessage) error {
	param := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(sr.queueURL),
		ReceiptHandle: aws.String(msg.msgReceiptHandle),
	}

	_, err := sr.sqs.DeleteMessageWithContext(sr.ctx, param)
	if err != nil {
		return errors.Wrapf(err, "cannot delete message with params. msg: %+v, params: %+v", msg, param)
	}

	return nil
}
