package myrpc

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// SenderConf contains info about config of message sender
type SenderConf struct {
	Queue QueueConf `yaml:"queue"`
}

// ReceiverConf contains info about config of message receiver
type ReceiverConf struct {
	Queue             QueueConf `yaml:"queue"`
	NumMsgsPerReceive int64     `yaml:"number_messages_per_receive"`
	VisibilityTimeout int64     `yaml:"visibility_timeout"`
	WaitTimeSeconds   int64     `yaml:"wait_time_seconds"`
}

// QueueConf contains info about message queue
type QueueConf struct {
	QueueRegion        string `yaml:"queue_region"`
	QueueBaseURL       string `yaml:"queue_base_url"`
	QueueName          string `yaml:"queue_name"`
	AWSAccessKeyID     string `yaml:"aws_access_key_id"`
	AWSSecretAccessKey string `yaml:"aws_secret_access_key"`
	SessionToken       string `yaml:"sqs_session_token"`
}

// ReceiverConfFromYamlFile returns ReceiverConf from given yaml conf file
func ReceiverConfFromYamlFile(filePath string) (*ReceiverConf, error) {
	var rc ReceiverConf
	bytes, err := fileToBytes(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid receiver conf file: %s", filePath)
	}

	if err := yaml.Unmarshal(bytes, &rc); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal receiver conf file: %s", filePath)
	}

	return &rc, nil
}

// SenderConfFromYamlFile returns ReceiverConf from given yaml conf file
func SenderConfFromYamlFile(filePath string) (*SenderConf, error) {
	var sc SenderConf
	bytes, err := fileToBytes(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid sender conf file: %s", filePath)
	}

	if err := yaml.Unmarshal(bytes, &sc); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal sender conf file: %s", filePath)
	}

	return &sc, nil
}

func fileToBytes(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open file: %s", filePath)
	}
	defer f.Close()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read file: %s", filePath)
	}

	return bytes, nil
}
