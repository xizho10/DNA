package main

import (
	"DNA/account"
	"DNA/common"
	"DNA/common/log"
	"DNA/core/code"
	"DNA/core/contract"
	"DNA/core/ledger"
	"DNA/core/store/ChainStore"
	"DNA/core/transaction"
	"DNA/crypto"
	"DNA/smartcontract"
	"DNA/smartcontract/service"
	"DNA/smartcontract/types"
	"DNA/vm/avm"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println("======")
	log.Init(log.Path, log.Stdout)
	log.Init(log.Path, log.Stdout)
	ledger.DefaultLedger = new(ledger.Ledger)
	ledger.DefaultLedger.Store = ChainStore.NewLedgerStore()
	ledger.DefaultLedger.Store.InitLedgerStore(ledger.DefaultLedger)
	transaction.TxStore = ledger.DefaultLedger.Store

	crypto.SetAlg("P256R1")
	type Person struct {
		Age    int    `json:Age`
		Gender string `json:Gender`
		Cge    int    `json:Cge`
	}
	p := &Person{22, "abc", 11}
	var data []byte
	if bs, err := json.Marshal(p); err == nil {
		data = bs
		fmt.Println(bs)
	}

	//return
	c, _ := common.HexToBytes("55c56b6c766b00527ac4616153c5760139007cc47603616263517cc4765b527cc46c766b51527ac46c766b51c300c301337e6c766b52527ac45cc57600017bc476510622416765223ac476526c766b52c3c47653012cc47654092247656e646572223ac476550122c476566c766b51c351c3c476570122c47658012cc476590622436765223ac4765a0131c4765b017dc46c766b53527ac46c766b53c3616516006c766b54527ac46203006c766b54c3616c756656c56b6c766b00527ac461006c766b51527ac4616c766b00c36c766b52527ac4006c766b53527ac46237006c766b52c36c766b53c3c36c766b54527ac4616c766b51c36c766b54c37e6c766b51527ac4616c766b53c351936c766b53527ac46c766b53c36c766b52c3c09f63c0ff6c766b51c36c766b55527ac46203006c766b55c3616c7566")
	//c, _ := common.HexToBytes("746b4c0400000000614c04e8030000744c0400000000948c6c766b947275620300744c0400000000948c6c766b947961748c6c766b946d6c7566")
	//c, _ := common.HexToBytes("746b4c0400000000614c04e8030000744c0400000000948c6c766b947275620300744c0400000000948c6c766b947961748c6c766b946d746c768c6b946d746c768c6b946d746c768c6b946d6c7566746b4c040000000061744c0400000000936c766b9479744c0401000000936c766b947993744c0400000000948c6c766b947275620300744c0400000000948c6c766b947961748c6c766b946d746c768c6b946d746c768c6b946d6c7566")
	engine := avm.NewExecutionEngine(nil, new(avm.ECDsaCrypto), nil, service.NewStateMachine(nil), common.Fixed64(0))
	engine.LoadCode(c, false)
	engine.Execute()
	val := engine.GetEvaluationStack().Pop().GetStackItem().GetByteArray()
	fmt.Println("Result:", string(val))
	//val := engine.GetEvaluationStack().Pop().GetStackItem().GetArray()
	//fmt.Println()
	//fmt.Println("LEN:", len(val),val)
	//for k, v := range val {
	//	fmt.Println(k, ": ", v.GetByteArray())
	//}
	//fmt.Println("{:", []byte(`{`))
	//fmt.Println("Age:", []byte(`"Age":`))
	//fmt.Println("Gender:", []byte(`,"Gender":`))
	//fmt.Println("Cge:", []byte(`,"Cge":`))
	//fmt.Println("}:", []byte(`}`))
	//fmt.Println(data)
	if false {
		fmt.Println(string(data))
	}

	//fmt.Println("Result:",val[0].GetByteArray(),val[1].GetByteArray())
	return
	fc := &code.FunctionCode{
		Code:           c,
		ParameterTypes: []contract.ContractParameterType{contract.Integer, contract.Integer},
		ReturnType:     contract.ContractParameterType(contract.Integer),
	}

	fmt.Println("CodeHash:", fc.CodeHash())
	dbCache := ChainStore.NewDBCache(ledger.DefaultLedger.Store.(*ChainStore.ChainStore))

	smartContract, err := smartcontract.NewSmartContract(&smartcontract.Context{
		Language:     0,
		StateMachine: service.NewStateMachine(dbCache),
		ReturnType:   contract.ContractParameterType(contract.Integer),
		Input:        c,
		CodeHash:     fc.CodeHash(),
	})
	smartContract.CodeHash = fc.CodeHash()
	if err != nil {

	}
	//	smartContract.DeployContractZX()

}
func makeDeployContractTransaction(codeStr string, language int) (string, error) {
	c, _ := common.HexToBytes(codeStr)
	fc := &code.FunctionCode{
		Code:           c,
		ParameterTypes: []contract.ContractParameterType{contract.Integer, contract.Integer},
		ReturnType:     contract.ContractParameterType(contract.Integer),
	}
	fc.CodeHash()
	acc, err := account.NewAccount()
	tx, err := transaction.NewDeployTransaction(fc, acc.ProgramHash, "test", "1.0", "user", "user@163.com", "test uint", types.LangType(byte(language)))
	if err != nil {
		return "Deploy smartcontract fail!", err
	}

	var buffer bytes.Buffer
	if err := tx.Serialize(&buffer); err != nil {
		fmt.Println("serialize registtransaction failed")
		return "", err
	}
	return hex.EncodeToString(buffer.Bytes()), nil
}

func makeInvokeTransaction(paramsStr, codeHashStr string) (string, error) {
	p, _ := common.HexToBytes(paramsStr)
	hash, _ := common.HexToBytes(codeHashStr)
	codeHash := common.BytesToUint160(hash)
	tx, err := transaction.NewInvokeTransaction(p, codeHash)
	if err != nil {
		return "Invoke smartcontract fail!", err
	}

	var buffer bytes.Buffer
	if err := tx.Serialize(&buffer); err != nil {
		fmt.Println("serialize registtransaction failed")
		return "", err
	}
	return hex.EncodeToString(buffer.Bytes()), nil
}
