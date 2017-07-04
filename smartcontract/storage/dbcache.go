package storage

import (
	"DNA/core/store"
	"DNA/smartcontract/states"
	"DNA/common"
	"math/big"
)

type DBCache interface {
	GetOrAdd(prefix store.DataEntryPrefix, key string, value states.IStateValueInterface) (states.IStateValueInterface, error)
	TryGet(prefix store.DataEntryPrefix, key string) (states.IStateValueInterface, error)
	GetWriteSet() *RWSet
	GetState(codeHash common.Uint160, loc common.Hash) (common.Hash, error)
	SetState(codeHash common.Uint160, loc, value common.Hash)
	GetCode(codeHash common.Uint160) ([]byte, error)
	SetCode(codeHash common.Uint160, code []byte)
	GetBalance(common.Uint160) *big.Int
	GetCodeSize(common.Uint160) int
	AddBalance(common.Uint160, *big.Int)
	Suicide(codeHash common.Uint160) bool
}
