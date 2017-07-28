package smartcontract

import (
	"DNA/common"
	sig "DNA/core/signature"
	"DNA/smartcontract/service"
	"DNA/smartcontract/storage"
	"DNA/smartcontract/types"
	"DNA/vm/avm"
	"DNA/vm/avm/interfaces"
	"math/big"
	//"DNA/vm/evm"
	"DNA/core/contract"
	"DNA/errors"
	//"DNA/vm/evm/abi"
	"DNA/common/log"
	"DNA/common/serialization"
	"DNA/core/asset"
	"DNA/core/ledger"
	"DNA/core/transaction"
	"DNA/smartcontract/states"
	"bytes"
	"fmt"
	"strconv"
	//"github.com/ethereum/go-ethereum/accounts/abi"
	//"github.com/ethereum/go-ethereum/accounts/abi"
	//"github.com/ethereum/go-ethereum/accounts/abi"
)

type SmartContract struct {
	Engine         Engine
	Code           []byte
	Input          []byte
	ParameterTypes []contract.ContractParameterType
	//ABI            abi.ABI
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
	ParameterTypes []contract.ContractParameterType
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
			//e = evm.NewExecutionEngine(context.DBCache, context.Time, context.BlockNumber, context.Gas)
		}

		return &SmartContract{
			Engine:         e,
			Code:           context.Code,
			CodeHash:       context.CodeHash,
			Input:          context.Input,
			Caller:         context.Caller,
			VMType:         vmType,
			ReturnType:     context.ReturnType,
			ParameterTypes: context.ParameterTypes,
		}, nil
	} else {
		return nil, errors.NewDetailErr(errors.NewErr("Not Support Language Type!"), errors.ErrNoCode, "")
	}

}

func (sc *SmartContract) DeployContract() ([]byte, error) {
	return sc.Engine.Create(sc.Caller, sc.Code)
}

func (sc *SmartContract) InvokeContract() (interface{}, error) {
	input, err := sc.InvokeParamsTransform()
	if err != nil {
		return nil, err
	}
	log.Error("==========input=========", input)
	sc.Engine.Call(sc.Caller, sc.CodeHash, input)
	return sc.InvokeResult()
}

func (sc *SmartContract) InvokeResult() (interface{}, error) {
	switch sc.VMType {
	case types.AVM:
		engine := sc.Engine.(*avm.ExecutionEngine)
		log.Error("==========type========", sc.ReturnType)
		log.Error("==========type========", engine.GetEvaluationStackCount())
		if engine.GetEvaluationStackCount() > 0 && avm.Peek(engine).GetStackItem() != nil {
			switch sc.ReturnType {
			case contract.Boolean:
				log.Error("=========Result==========", avm.Peek(engine))
				return avm.PopBoolean(engine), nil
			case contract.Integer:
				return avm.PopInt(engine), nil
			case contract.ByteArray:
				log.Error("=========Result ByteArray==========", string(avm.Peek(engine).GetStackItem().GetByteArray()))
				return string(avm.PopByteArray(engine)), nil
			case contract.Hash160, contract.Hash256:
				return common.ToHexString(common.ToArrayReverse(avm.PopByteArray(engine))), nil
			case contract.Object:
				log.Error("==============Object============", avm.Peek(engine).GetStackItem())
				data := avm.PopInteropInterface(engine)
				switch data.(type) {
				case *ledger.Header:
					log.Error("==============ledger.Header============", data.(*ledger.Header))
					return service.GetHeaderInfo(data.(*ledger.Header)), nil
				case *ledger.Block:
					log.Error("==============ledger.Block============", data.(*ledger.Block))
					return service.GetBlockInfo(data.(*ledger.Block)), nil
				case *transaction.Transaction:
					log.Error("==============transaction.Transaction============", data.(*transaction.Transaction))
					return service.GetTransactionInfo(data.(*transaction.Transaction)), nil
				case *states.AccountState:
					return service.GetAccountInfo(data.(*states.AccountState)), nil
				case *asset.Asset:
					return service.GetAssetInfo(data.(*asset.Asset)), nil
				default:
					return data, nil
				}
			}
		}
	case types.EVM:
	}
	return nil, nil
}

func (sc *SmartContract) InvokeParamsTransform() ([]byte, error) {
	fmt.Println("===========InvokeParamsTransform=============")
	switch sc.VMType {
	case types.AVM:
		builder := avm.NewParamsBuilder(new(bytes.Buffer))
		fmt.Println("==========sc.Input=============", sc.Input)
		b := bytes.NewBuffer(sc.Input)
		for _, k := range sc.ParameterTypes {
			switch k {
			case contract.Boolean:
				p, err := serialization.ReadBool(b)
				if err != nil {
					return nil, err
				}
				builder.EmitPushBool(p)
			case contract.Integer:
				p, err := serialization.ReadVarBytes(b)
				if err != nil {
					return nil, err
				}
				i, err := strconv.ParseInt(string(p), 10, 64)
				if err != nil {
					return nil, err
				}
				fmt.Println("===========p=============", int64(i))
				builder.EmitPushInteger(int64(i))
			case contract.ByteArray:
				p, err := serialization.ReadVarBytes(b)
				if err != nil {
					return nil, err
				}
				builder.EmitPushByteArray(p)
			}
		}
		builder.EmitPushCall(sc.CodeHash.ToArray())
		fmt.Println("=================builder======================", builder.ToArray())
		return builder.ToArray(), nil
	case types.EVM:
	}
	return nil, nil
}
