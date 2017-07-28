package types

import (
	"math/big"
	"DNA/vm/avm/interfaces"
)

type ByteArray struct {
	value []byte
}

func NewByteArray(value []byte) *ByteArray {
	var ba ByteArray
	ba.value = value
	return &ba
}

func (ba *ByteArray) Equals(other StackItemInterface) bool {
	if _, ok := other.(*ByteArray); !ok {
		return false
	}
	a1 := ba.value
	a2 := other.GetByteArray()
	l1 := len(a1)
	l2 := len(a2)
	if l1 != l2 {
		return false
	}
	for i:=0; i<l1; i++ {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func (ba *ByteArray) GetBigInteger() *big.Int {
	bi := new(big.Int)
	return bi.SetBytes(ToArrayReverse(ba.value))
}

func (ba *ByteArray) GetBoolean() bool{
	for _, b := range ba.value {
		if b != 0 {
			return true
		}
	}
	return false
}

func (ba *ByteArray) GetByteArray() []byte {
	return ba.value
}

func (ba *ByteArray) GetInterface() interfaces.IInteropInterface {
	return nil
}

func (ba *ByteArray) GetArray() []StackItemInterface {
	return []StackItemInterface{ba}
}

func ToArrayReverse(arr []byte) []byte {
	l := len(arr)
	var x []byte = make([]byte, l)
	for i, j := 0, l-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = arr[j], arr[i]
	}
	return x
}
