package myrpc

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

func newSQSClient(ctx context.Context, conf QueueConf) (client *sqs.SQS, queueURL string, err error) {
	sqsSession, err := initSQSSession(conf)
	if err != nil {
		return nil, "", errors.Wrapf(err, "cannot init sqs session with config: %+v", conf)
	}

	client = sqs.New(sqsSession)
	queueInfo, err := client.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(conf.QueueName)})
	if err != nil {
		return nil, "", errors.Wrapf(err, "cannot get queue url %s", conf.QueueName)
	}

	return client, *queueInfo.QueueUrl, nil
}

func initSQSSession(conf QueueConf) (*session.Session, error) {
	cred := credentials.NewStaticCredentials(conf.AWSAccessKeyID, conf.AWSSecretAccessKey, conf.SessionToken)
	return session.New(&aws.Config{
		Credentials: cred,
		Region:      aws.String(conf.QueueRegion),
		Endpoint:    aws.String(conf.QueueBaseURL),
	}), nil
}
