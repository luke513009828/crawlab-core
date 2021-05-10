package errors

func NewNodeError(msg string) (err error) {
	return NewError(ErrorPrefixNode, msg)
}

var ErrorNodeUnregistered = NewNodeError("unregistered")
var ErrorNodeServiceNotExists = NewNodeError("service not exists")
var ErrorNodeInvalidType = NewNodeError("invalid type")
var ErrorNodeInvalidStatus = NewNodeError("invalid status")
var ErrorNodeInvalidCode = NewNodeError("invalid code")
var ErrorNodeInvalidNodeKey = NewNodeError("invalid node key")
var ErrorNodeStreamNotFound = NewNodeError("stream not found")
var ErrorNodeMonitorError = NewNodeError("monitor error")