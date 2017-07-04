package smartcontract

import (
	"DNA/common"
	"math/big"
	"DNA/vm/avm/interfaces"
	sig "DNA/core/signature"
	"DNA/smartcontract/storage"
	"DNA/smartcontract/service"
	"DNA/smartcontract/types"
	"DNA/vm/avm"
	"DNA/vm/evm"
	"DNA/errors"
	"DNA/core/contract"
)

type SmartContract struct {
	Engine     Engine
	Code       []byte
	Input      []byte
	Caller     common.Uint160
	CodeHash   common.Uint160
	VMType     types.VmType
	ReturnType contract.ContractParameterType
}

type Context struct {
	Language       types.LangType
	Caller         common.Uint160
	StateMachine   *service.StateMachine
	DBCache        storage.DBCache
	Code           []byte
	Input          []byte
	CodeHash       common.Uint160
	Time           *big.Int
	BlockNumber    *big.Int
	CacheCodeTable interfaces.ICodeTable
	SignableData   sig.SignableData
	Gas            common.Fixed64
	ReturnType     contract.ContractParameterType
}

type Engine interface {
	Create(caller common.Uint160, code []byte) ([]byte, error)
	Call(caller common.Uint160, codeHash common.Uint160, input []byte) ([]byte, error)
}

func NewSmartContract(context *Context) (*SmartContract, error) {
	if vmType, ok := types.LangVm[context.Language]; ok {
		var e Engine
		switch vmType {
		case types.AVM:
			e = avm.NewExecutionEngine(
				context.SignableData,
				new(avm.ECDsaCrypto),
				context.CacheCodeTable,
				context.StateMachine,
				context.Gas,
			)
		case types.EVM:
			e = evm.NewExecutionEngine(context.DBCache, context.Time, context.BlockNumber, context.Gas)
		}

		return &SmartContract{
			Engine: e,
			Code: context.Code,
			CodeHash: context.CodeHash,
			Input: context.Input,
			Caller: context.Caller,
			VMType: vmType,
			ReturnType: context.ReturnType,
		}, nil
	} else {
		return nil, errors.NewDetailErr(errors.NewErr("Not Support Language Type!"), errors.ErrNoCode, "")
	}

}

func (sc *SmartContract) DeployContract() ([]byte, error) {
	return sc.Engine.Create(sc.Caller, sc.Code)
}

func (sc *SmartContract) InvokeContract() (interface{}, error) {
	sc.Engine.Call(sc.Caller, sc.CodeHash, sc.Input)
	return sc.InvokeResult()
}

func (sc *SmartContract) InvokeResult() (interface{}, error) {
	switch sc.VMType {
	case types.AVM:
		engine := sc.Engine.(*avm.ExecutionEngine)
		if engine.GetEvaluationStack().Count() > 0 {
			switch sc.ReturnType {
			case contract.Boolean:
				return avm.PopBoolean(engine), nil
			case contract.Integer:
				return avm.PopInt(engine), nil
			case contract.ByteArray:
				return avm.PopByteArray(engine), nil
			}
		}
	case types.EVM:
	}
	return nil, nil
}