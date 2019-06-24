package myrpc

import "encoding/json"

// RPCMessage represents message of this RPC
type RPCMessage struct {
	SvrName ServiceName `json:"service_name"`
	MthName MethodName  `json:"method_name"`
	Payload []byte      `json:"payload"`
	// use for delete message
	msgReceiptHandle string
}

// ToJSON converts RPCMessage to json in string format
func (msg *RPCMessage) ToJSON() (string, error) {
	bytes, err := json.Marshal(msg)
	return string(bytes), err
}

// JSONToRPCMsg converts json in string format to RPCMessage
func JSONToRPCMsg(jsonStr string) (*RPCMessage, error) {
	var msg RPCMessage
	if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
