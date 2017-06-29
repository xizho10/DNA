package evm

import (
	"DNA/common"
	. "DNA/vm/evm/common"
	"math/big"
)

type StateDB interface {
	GetState(common.Uint160, Hash) Hash
	SetState(common.Uint160, Hash, Hash)

	GetCode(common.Uint160) []byte
	SetCode(common.Uint160, []byte)
	GetCodeSize(common.Uint160) int

	GetBalance(common.Uint160) *big.Int
	AddBalance(common.Uint160, *big.Int)

	Suicide(common.Uint160) bool
}
