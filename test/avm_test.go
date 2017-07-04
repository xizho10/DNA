package test

import (
	"testing"
	"DNA/vm/avm"
	"fmt"
	"DNA/core/ledger"
	"DNA/core/store/ChainStore"
	"DNA/core/transaction"
	"DNA/common"
)

func TestAVM(t *testing.T) {
	fmt.Println("//**************************************************************************")
	code := `00C56B51616C7566`
	ledger.DefaultLedger = new(ledger.Ledger)
	ledger.DefaultLedger.Store = ChainStore.NewLedgerStore()
	ledger.DefaultLedger.Store.InitLedgerStore(ledger.DefaultLedger)
	transaction.TxStore = ledger.DefaultLedger.Store
	//crypto.SetAlg(""+crypto.P256R1)
	//var cryptos interfaces.ICrypto
	//cryptos = new(avm.ECDsaCrypto)
	avm := avm.NewExecutionEngine(nil, nil, nil, nil, common.Fixed64(0))
	codes, _ := common.HexToBytes(code)
	fmt.Println("codes:", codes)
	_, err := avm.Create(common.Uint160{}, codes)
	fmt.Println("ret:", avm.GetEvaluationStack().Pop().GetStackItem().GetByteArray(), "err:", err)
}
