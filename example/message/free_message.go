package message

// FreeMessageIn is a example about payload with non-protobuf message
type FreeMessageIn struct {
	Msg string `json:"msg"`
}

// FreeMessageOut is a example about payload with non-protobuf message
type FreeMessageOut struct {
	Msg string `json:"msg"`
}
