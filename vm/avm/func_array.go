package avm

import (
	. "DNA/vm/avm/errors"
	"fmt"
	"DNA/vm/avm/types"
	"reflect"
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
	if index < 0 {
		return FAULT, ErrFault
	}
	items := PopArray(e)
	if reflect.TypeOf(items).Kind() != reflect.Slice || reflect.TypeOf(items).Kind() != reflect.Array {
		return FAULT, ErrNotArray
	}
	if index >= len(items) {
		return FAULT, ErrOverLen
	}
	PushData(e, items[index])
	return NONE, nil
}

func opSetItem(e *ExecutionEngine) (VMState, error) {
	newItem := Pop(e)
	index := PopInt(e)
	arrItem := PopStackItem(e)
	if arrItem == nil || reflect.TypeOf(arrItem).Kind() != reflect.Slice || reflect.TypeOf(arrItem).Kind() != reflect.Array {
		return NONE, nil
	}
	items := arrItem.GetArray()
	if index < 0 || index >= len(items) {
		return FAULT, fmt.Errorf("%v", "set item invalid index")
	}
	items[index] = newItem.GetStackItem()
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


