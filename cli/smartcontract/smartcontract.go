package smartcontract

import (
	"DNA/core/transaction"
	"DNA/core/code"
	"DNA/core/contract"
	"DNA/common"
	httpjsonrpc "DNA/net/httpjsonrpc"
	"DNA/smartcontract/types"
	"github.com/urfave/cli"
	. "DNA/cli/common"
	"fmt"
	"os"
	"bytes"
	"encoding/hex"
	"DNA/account"
)

func makeDeployContractTransaction(codeStr string, language int) (string, error) {
	c, _ := common.HexToBytes(codeStr)
	fc := &code.FunctionCode{
		Code: c,
		ParameterTypes: []contract.ContractParameterType{contract.Integer, contract.Integer},
		ReturnType: contract.ContractParameterType(contract.Integer),
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

func contractAction(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	var err error
	var txHex string
	deploy := c.Bool("deploy")
	invoke := c.Bool("invoke")
	if !deploy && !invoke {
		fmt.Println("missing --deploy -d or --invoke -i")
		return nil
	}

	if deploy {
		codeStr := c.String("code")
		language := c.Int("language")
		if codeStr == "" {
			fmt.Println("missing args [--code] or [--language]")
			return nil
		}
		txHex, err = makeDeployContractTransaction(codeStr, language)
		if err != nil {
			fmt.Println(err)
		}
	}
	if invoke {
		paramsStr := c.String("params")
		codeHashStr := c.String("codeHash")
		txHex, err = makeInvokeTransaction(paramsStr, codeHashStr)
		if err != nil {
			fmt.Println(err)
		}
	}
	resp, err := httpjsonrpc.Call(Address(), "sendrawtransaction", 0, []interface{}{txHex})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	FormatOutput(resp)
	return nil
}


func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "contract",
		Usage:       "deploy or invoke your smartcontract ",
		Description: "you could deploy or invoke your smartcontract.",
		ArgsUsage:   "[args]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "deploy, d",
				Usage: "deploy smartcontract",
			},
			cli.BoolFlag{
				Name:  "invoke, i",
				Usage: "invoke smartcontract",
			},
			cli.StringFlag{
				Name:  "code, c",
				Usage: "deploy contract code",
			},
			cli.IntFlag{
				Name:  "language, l",
				Usage: "deploy contract compiler contract language",
			},
			cli.StringFlag{
				Name:  "params, p",
				Usage: "invoke contract compiler contract params",
			},
			cli.StringFlag{
				Name:  "codeHash, a",
				Usage: "invoke contract compiler contract code hash",
			},
		},
		Action: contractAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			PrintError(c, err, "smartcontract")
			return cli.NewExitError("", 1)
		},
	}
}

