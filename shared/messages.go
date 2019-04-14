package shared

type RPCMessage struct {
	SvName ServiceName `json:"service_name"`
	MtName MethodName `json:"method_name"`
	Payload []byte `json:"payload"`
}

type GetUserRequest struct {
	UserID uint64 `json:"user_id"`
}

type GetUserResponse struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
}
