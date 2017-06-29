package smartcontract

import (
	"DNA/vm/avm"
	"DNA/vm/evm"
	"github.com/pkg/errors"
	"DNA/common"
	"DNA/smartcontract/service"
)

type LangType byte

const (
	CSharp  LangType = iota
	Solidity
)

type VmType byte

const (
	AVM VmType = iota
	EVM
)

var LangVm map[LangType]VmType

func init() {
	LangVm = make(map[LangType]VmType, 0)
	LangVm[CSharp] = AVM
	LangVm[Solidity] = EVM
}

type SmartContract struct {
	Engine Engine
	Code []byte
	Input []byte
	CodeHash common.Uint160
}

func (sc *SmartContract) DeployContract() ([]byte, error) {
	sc.Engine.Call(sc.Code, sc.Input)
	return nil, nil
}

func (sc *SmartContract) InvokeContract() ([]byte, error) {
	sc.Engine.Call(sc.CodeHash, sc.Input)
	return nil, nil
}