package states

import (
	"io"
	. "DNA/errors"
	"DNA/core/code"
	"bytes"
)

type ContractState struct {
	Code *code.FunctionCode
	Name string
	Version string
	Author string
	Email string
	Description string
	*StateBase
}

func(contractState *ContractState)Serialize(w io.Writer) error {
	contractState.Code.Serialize(w)
	return nil
}

func(contractState *ContractState)Deserialize(r io.Reader) error {
	f := new(code.FunctionCode)
	err := f.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "ContractState Code Deserialize fail.")
	}
	contractState.Code = f

	return nil
}

func(contractState *ContractState) ToArray() []byte {
	b := new(bytes.Buffer)
	contractState.Serialize(b)
	return b.Bytes()
}


