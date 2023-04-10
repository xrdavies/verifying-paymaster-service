package errors

var (
	REJECTED_BY_TYPE = -32500
)

type RPCError struct {
	code    int
	message string
	data    any
}

func NewRPCError(code int, message string, data any) error {
	return &RPCError{code, message, data}
}

func (e *RPCError) Error() string {
	return e.message
}

func (e *RPCError) Data() any {
	return e.data
}

func (e *RPCError) Code() int {
	return e.code
}
