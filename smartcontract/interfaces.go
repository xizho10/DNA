package smartcontract

import "DNA/common"

type CallContract interface {
	DeployContract() ([]byte, error)
	InvokeContract() ([]byte, error)
}

type Engine interface {
	Create(code []byte, input []byte) ([]byte, error)
	Call(codeHash, input []byte) ([]byte, error)
}
