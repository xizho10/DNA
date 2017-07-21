package avm

func opCat(e *ExecutionEngine) (VMState, error) {
	b2 :=PopByteArray(e)
	b1 :=PopByteArray(e)

	r := ByteArrZip(b1, b2, CAT)
	PushData(e, r)
	return NONE, nil
}

func opSubStr(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	if count < 0 { return FAULT, nil }
	index := PopInt(e)
	if index < 0 { return FAULT, nil }
	s  :=PopByteArray(e)
	l1 := index + count
	l2 := len(s)
	if l1 > l2 {
		return FAULT, nil
	}
	b := s[index : l2-l1+1]
	PushData(e, b)
	return NONE, nil
}

func opLeft(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	if count < 0 { return FAULT, nil }
	s := PopByteArray(e)
	b := s[:count]
	PushData(e, b)
	return NONE, nil
}

func opRight(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	if count < 0 { return FAULT, nil }
	s := PopByteArray(e)
	l := len(s)
	if count > l { return FAULT, nil }
	b := s[l-count:]
	PushData(e, b)
	return NONE, nil
}

func opSize(e *ExecutionEngine) (VMState, error) {
	x := Peek(e).GetStackItem()
	PushData(e, len(x.GetByteArray()))
	return NONE, nil
}
