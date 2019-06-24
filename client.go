package myrpc

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
)

// MessageSender is the interface to send a message service
type MessageSender interface {
	// SendAsyncMsg is called in case of publish/subscribe pattern
	SendAsyncMsg(msg *RPCMessage) error
	// SendSyncMsg is called in case of request/reply pattern
	SendSyncMsg(in *RPCMessage, out interface{}) error
}

// PayloadEncodeFnc encodes data to bytes
type PayloadEncodeFnc func(data interface{}) ([]byte, error)

// RPCClient represents client of this RPC
type RPCClient struct {
	sender        MessageSender
	ctx           context.Context
	payloadEncode PayloadEncodeFnc
}

// NewRPCClient returns new client from config
func NewRPCClient(ctx context.Context, sender MessageSender) *RPCClient {
	return &RPCClient{
		sender:        sender,
		ctx:           ctx,
		payloadEncode: json.Marshal,
	}
}

// ReplacePayloadEncoder replaces payload encode function of rpc client
func (c *RPCClient) ReplacePayloadEncoder(encFnc PayloadEncodeFnc) {
	c.payloadEncode = encFnc
}

// SendAsyncMsg sends message to message service asynchronously,
// that means no waiting response from server.
// If no encodeFnc given, use client default encode instead.
func (c *RPCClient) SendAsyncMsg(svr ServiceName, mth MethodName, in interface{}, encodeFnc PayloadEncodeFnc) error {
	if encodeFnc == nil {
		// fallback to client default encode func
		encodeFnc = c.payloadEncode
	}

	payload, err := encodeFnc(in)
	if err != nil {
		return errors.Wrapf(err, "cannot encode payload: %+v", in)
	}

	rpcMsg := RPCMessage{
		SvrName: svr,
		MthName: mth,
		Payload: payload,
	}

	return c.sender.SendAsyncMsg(&rpcMsg)
}

// SendSyncMsg sends message to message service synchronously,
// that means it is blocked until received response from server
func (c *RPCClient) SendSyncMsg(svr ServiceName, mth MethodName, in interface{}, out interface{}, encodeFnc PayloadEncodeFnc) error {
	if encodeFnc == nil {
		encodeFnc = c.payloadEncode
	}

	payload, err := encodeFnc(in)
	if err != nil {
		return errors.Wrapf(err, "cannot encode payload: %+v", in)
	}

	rpcMsg := RPCMessage{
		SvrName: svr,
		MthName: mth,
		Payload: payload,
	}

	return c.sender.SendSyncMsg(&rpcMsg, out)
}
