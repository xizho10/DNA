package errors


import "errors"

var (
	ErrBadValue           = errors.New("bad value")
	ErrBadType            = errors.New("bad type")
	ErrOverLen	      = errors.New("the count over the size")
	ErrLittleLen          = errors.New("the count too little")
	ErrFault	      = errors.New("the exeution meet fault")
	ErrNotSupportService  = errors.New("the service is not registered")
	ErrNotSupportOpCode   = errors.New("does not support the operation code")
	ErrOverLimitStack     = errors.New("the stack over max size")
	ErrOverMaxItemSize    = errors.New("the item over max size")
	ErrOutOfGas           = errors.New("out of gas")
	ErrNotArray           = errors.New("not array")
)
