package service

import (
	"DNA/common"
	"bytes"
)

type StorageContext struct {
	codeHash *common.Uint160
}

func NewStorageContext(codeHash *common.Uint160) *StorageContext {
	var storageContext StorageContext
	storageContext.codeHash = codeHash
	return &storageContext
}

func (sc *StorageContext) ToArray() ([]byte) {
	b := new(bytes.Buffer)
	sc.codeHash.Serialize(b)
	return b.Bytes()
}
