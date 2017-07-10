package payload

import (
	"DNA/common/serialization"
	. "DNA/errors"
	"errors"
	"io"
)

type StateUpdate struct {
	Namespace []byte
	Key       []byte
	Value     []byte
}

func (su *StateUpdate) Data() []byte {
	return []byte{0}
}

func (su *StateUpdate) Serialize(w io.Writer) error {
	err := serialization.WriteVarBytes(w, su.Namespace)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[StateUpdate], Namespace serialize failed.")
	}

	err = serialization.WriteVarBytes(w, su.Key)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[StateUpdate], key serialize failed.")
	}

	err = serialization.WriteVarBytes(w, su.Value)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[StateUpdate], value serialize failed.")
	}

	return nil
}

func (su *StateUpdate) Deserialize(r io.Reader) error {
	var err error

	su.Namespace, err = serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(errors.New("[StateUpdate], Namespace deserialize failed."), ErrNoCode, "")
	}

	su.Key, err = serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(errors.New("[StateUpdate], key deserialize failed."), ErrNoCode, "")
	}

	su.Value, err = serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(errors.New("[StateUpdate], value deserialize failed."), ErrNoCode, "")
	}

	//su.Updater = *new(common.Uint160)
	//err = su.Updater.Deserialize(r)
	//if err != nil {
	//	return NewDetailErr(err, ErrNoCode, "[StateUpdate], updater Deserialize failed.")
	//}

	return nil
}
