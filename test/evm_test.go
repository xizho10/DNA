package test

import (
	"testing"
	"strings"
	"fmt"
	"DNA/common"
	"DNA/vm/evm/abi"
	"DNA/vm/evm"
	"math/big"
	"DNA/crypto"
	client "DNA/account"
)

const (
	ABI = `[{"constant":false,"inputs":[{"name":"x","type":"uint256"},{"name":"y","type":"uint256"},{"name":"z","type":"uint256"}],"name":"set","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"}]`
	BIN = `6060604052341561000c57fe5b5b60f08061001b6000396000f30060606040526000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806343b0e8df1460445780636d4ce63c146073575bfe5b3415604b57fe5b607160048080359060200190919080359060200190919080359060200190919050506096565b005b3415607a57fe5b608060b1565b6040518082815260200191505060405180910390f35b8260008190555081600181905550806002819055505b505050565b6000600254600154600054010290505b905600a165627a7a723058202607fc4aedbd9b26fa8cb43a115002b840ae3bc3d9c054ce55cc0b6e748c4a890029`
)
func TestEvm(t *testing.T) {
	parsed, err := abi.JSON(strings.NewReader(ABI))
	if err != nil {
		t.Fatal("parsed error:", err)
	}
	fmt.Println(parsed)
	input, err := parsed.Pack("")
	if err != nil {
		t.Fatal("input error:", err)
	}
	fmt.Println("input:", input)
	evm := evm.NewExecutionEngine(nil, big.NewInt(1), big.NewInt(1), common.Fixed64(0))
	code, _ := common.HexToBytes(BIN)
	crypto.SetAlg(fmt.Sprintf("%v", crypto.P256R1))
	account, _ := client.NewAccount()

	codeHash, _ := common.ToCodeHash(code)
	evm.Create(account.ProgramHash, code)

	input, err = parsed.Pack("set", 4, 2, 6)
	fmt.Println("input:", input)

	ret, err := evm.Call(account.ProgramHash, codeHash, input)
	fmt.Println("ret:", ret)

	input, err = parsed.Pack("get")

	fmt.Println("input:", input)

	ret, err = evm.Call(account.ProgramHash, codeHash, input)

	fmt.Println("ret:", ret)
}


