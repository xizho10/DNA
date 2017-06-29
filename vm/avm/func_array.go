package avm

import (
	. "DNA/vm/avm/errors"
	"fmt"
	"DNA/vm/avm/types"
)

func opArraySize(e *ExecutionEngine) (VMState, error) {
	arr := PopArray(e)
	PushData(e, len(arr))
	return NONE, nil
}

func opPack(e *ExecutionEngine) (VMState, error) {
	size := PopInt(e)
	if size > e.evaluationStack.Count() {
		return FAULT, ErrBadValue
	}
	items := NewStackItems()
	for i := 0; i< size; i++ {
		items = append(items, PopStackItem(e))
	}
	PushData(e, items)
	return NONE, nil
}

func opUnpack(e *ExecutionEngine) (VMState, error) {
	arr := PopArray(e)
	l := len(arr)
	for i := l - 1; i >= 0; i-- {
		Push(e, NewStackItem(arr[i]))
	}
	PushData(e, l)
	return NONE, nil
}

func opPickItem(e *ExecutionEngine) (VMState, error) {
	index := PopInt(e)
	items := PopArray(e)
	if index >= len(items) {
		return FAULT, ErrOverLen
	}
	PushData(e, items[index])
	return NONE, nil
}

func opSetItem(e *ExecutionEngine) (VMState, error) {
	newItem := Pop(e)
	index := PopInt(e)
	arrItem := PopArray(e)
	if index >= len(arrItem) {
		return FAULT, fmt.Errorf("%v", "set item invalid index")
	}
	arrItem[index] = newItem.GetStackItem()
	return NONE, nil
}

func opNewArray(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	items := NewStackItems();
	for i :=0; i <count; i++ {
		items = append(items, types.NewBoolean(false))
	}
	PushData(e, items)
	return NONE, nil
}


