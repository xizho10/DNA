package statedb

import (
	"DNA/smartcontract/storage"
	"DNA/common"
	. "DNA/vm/evm/common"
	"DNA/smartcontract/states"
	"DNA/core/store"
	"math/big"
)

type StateDB struct {
	RWSet *storage.RWSet
}

func NewStateDB() *StateDB {
	var stateDB StateDB
	stateDB.RWSet = storage.NewRWSet()
	return &stateDB
}

func (s *StateDB) GetState(codeHash common.Uint160, loc Hash) Hash {
	key := states.NewStorageKey(&codeHash, loc.Bytes())
	keyStr := storage.KeyToStr(key)
	value, _ := s.RWSet.TryGet(store.DataEntryPrefix(0), keyStr)
	if value == nil {
		return [32]byte{}
	}
	return BytesToHash(value.(*states.StorageItem).Value)
}

func (s *StateDB) SetState(codeHash common.Uint160, loc, value Hash) {
	key := states.NewStorageKey(&codeHash, loc.Bytes())
	item := states.NewStorageItem(value.Bytes())
	keyStr := storage.KeyToStr(key)
	s.RWSet.Put(keyStr, item)
}

func (s *StateDB) GetCode(codeHash common.Uint160) []byte {
	skey := storage.KeyToStr(&codeHash)
	value, err := s.RWSet.TryGet(store.ST_Contract, skey)
	if err != nil {
		return []byte{}
	}
	return value.(*states.StorageItem).Value
}

func (s *StateDB) SetCode(codeHash common.Uint160, code []byte) {
	item := states.NewStorageItem(code)
	skey := storage.KeyToStr(&codeHash)
	s.RWSet.Put(skey, item)

}

func (s *StateDB) GetBalance(common.Uint160) *big.Int {
	return big.NewInt(100)
}

func (s *StateDB) GetCodeSize(common.Uint160) int {
	return 0
}

func (s *StateDB) AddBalance(common.Uint160, *big.Int) {
}

func (s *StateDB) Suicide(codeHash common.Uint160) bool {
	skey := storage.KeyToStr(&codeHash)
	s.RWSet.Delete(skey)
	return true
}

