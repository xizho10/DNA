package service

import (
	"DNA/core/ledger"
	"DNA/common"
	"math/big"
	"errors"
	"fmt"
	"DNA/vm/avm/types"
	"DNA/core/transaction"
	"DNA/smartcontract/states"
	"DNA/vm/avm"
	"DNA/common/log"
)

type StateReader struct {
	serviceMap map[string]func(*avm.ExecutionEngine) (bool, error)
}

func NewStateReader() *StateReader {
	var stateReader StateReader
	stateReader.serviceMap = make(map[string]func(*avm.ExecutionEngine) (bool, error), 0)
	stateReader.Register("AntShares.Blockchain.GetHeight", stateReader.BlockChainGetHeight)
	stateReader.Register("AntShares.Blockchain.GetHeader", stateReader.BlockChainGetHeader)
	stateReader.Register("AntShares.Blockchain.GetBlock", stateReader.BlockChainGetBlock)
	stateReader.Register("AntShares.Blockchain.GetTransaction", stateReader.BlockChainGetTransaction)
	stateReader.Register("AntShares.Blockchain.GetAccount", stateReader.BlockChainGetAccount)
	stateReader.Register("AntShares.Blockchain.GetValidators", stateReader.BlockChainGetValidators)
	stateReader.Register("AntShares.Blockchain.GetAsset", stateReader.BlockChainGetAsset)

	stateReader.Register("AntShares.Header.GetHash", stateReader.HeaderGetHash);
	stateReader.Register("AntShares.Header.GetVersion", stateReader.HeaderGetVersion);
	stateReader.Register("AntShares.Header.GetPrevHash", stateReader.HeaderGetPrevHash);
	stateReader.Register("AntShares.Header.GetMerkleRoot", stateReader.HeaderGetMerkleRoot);
	stateReader.Register("AntShares.Header.GetTimestamp", stateReader.HeaderGetTimestamp);
	stateReader.Register("AntShares.Header.GetConsensusData", stateReader.HeaderGetConsensusData);
	stateReader.Register("AntShares.Header.GetNextConsensus", stateReader.HeaderGetNextConsensus);

	stateReader.Register("AntShares.Block.GetTransactionCount", stateReader.BlockGetTransactionCount);
	stateReader.Register("AntShares.Header.GetTransactions", stateReader.BlockGetTransactions);
	stateReader.Register("AntShares.Block.GetTransaction", stateReader.BlockGetTransaction);

	stateReader.Register("AntShares.Transaction.GetHash", stateReader.TransactionGetHash);
	stateReader.Register("AntShares.Transaction.GetType", stateReader.TransactionGetType);
	stateReader.Register("AntShares.Transaction.GetAttributes", stateReader.TransactionGetAttributes);
	stateReader.Register("AntShares.Transaction.GetInputs", stateReader.TransactionGetInputs);
	stateReader.Register("AntShares.Transaction.GetOutputs", stateReader.TransactionGetOutputs);
	stateReader.Register("AntShares.Transaction.GetReferences", stateReader.TransactionGetReferences);

	stateReader.Register("AntShares.Attribute.GetUsage", stateReader.AttributeGetUsage);
	stateReader.Register("AntShares.Attribute.GetData", stateReader.AttributeGetData);

	stateReader.Register("AntShares.Input.GetHash", stateReader.InputGetHash);
	stateReader.Register("AntShares.Input.GetIndex", stateReader.InputGetIndex);

	stateReader.Register("AntShares.Output.GetAssetId", stateReader.OutputGetAssetId);
	stateReader.Register("AntShares.Output.GetValue", stateReader.OutputGetValue);
	stateReader.Register("AntShares.Output.GetScriptHash", stateReader.OutputGetCodeHash);

	stateReader.Register("AntShares.Account.GetScriptHash", stateReader.AccountGetCodeHash);
	stateReader.Register("AntShares.Account.GetBalance", stateReader.AccountGetBalance);

	stateReader.Register("AntShares.Asset.GetAssetId", stateReader.AssetGetAssetId);
	stateReader.Register("AntShares.Asset.GetAssetType", stateReader.AssetGetAssetType);
	stateReader.Register("AntShares.Asset.GetAmount", stateReader.AssetGetAmount);
	stateReader.Register("AntShares.Asset.GetAvailable", stateReader.AssetGetAvailable);
	stateReader.Register("AntShares.Asset.GetPrecision", stateReader.AssetGetPrecision);
	stateReader.Register("AntShares.Asset.GetOwner", stateReader.AssetGetOwner);
	stateReader.Register("AntShares.Asset.GetAdmin", stateReader.AssetGetAdmin);
	stateReader.Register("AntShares.Asset.GetIssuer", stateReader.AssetGetIssuer);

	stateReader.Register("AntShares.Contract.GetScript", stateReader.ContractGetCode);

	stateReader.Register("AntShares.Storage.GetContext", stateReader.StorageGetContext);

	return &stateReader
}


func (s *StateReader) Register(methodName string, handler func(*avm.ExecutionEngine) (bool, error)) bool {
	if _, ok := s.serviceMap[methodName]; ok {
		return false
	}
	s.serviceMap[methodName] = handler
	return true
}

func (s *StateReader) GetServiceMap() map[string]func(*avm.ExecutionEngine) (bool, error) {
	return s.serviceMap
}

func (s *StateReader) BlockChainGetHeight(e *avm.ExecutionEngine) (bool, error) {
	var i uint32
	if ledger.DefaultLedger == nil {
		i = 0
	}else {
		log.Error("[BlockChainGetHeight] DefaultLedger Store:", i)
		i = ledger.DefaultLedger.Store.GetHeight()
	}
	log.Error("[BlockChainGetHeight] height:", i)
	avm.PushData(e, i)
	return true, nil
}

func (s *StateReader) BlockChainGetHeader(e *avm.ExecutionEngine) (bool, error) {
	data := avm.PopByteArray(e)
	var (
		header *ledger.Header
		err error
	)
	l := len(data)
	if l <= 5 {
		b := new(big.Int)
		height := uint32(b.SetBytes(data).Int64())
		if ledger.DefaultLedger != nil {
			hash, err := ledger.DefaultLedger.Store.GetBlockHash(height)
			if err != nil { return false, err }
			header, err = ledger.DefaultLedger.Store.GetHeader(hash)
			if err != nil { return false, err }
		}else {
			header = nil
		}
	}else if l == 32 {
		hash, _ := common.Uint256ParseFromBytes(data)
		if ledger.DefaultLedger != nil {
			header, err = ledger.DefaultLedger.Store.GetHeader(hash)
			if err != nil { return false, err }
		}else {
			header = nil
		}
	}else {
		return false, errors.New("The data length is error in function blockchaningetheader!")
	}
	log.Error("[BlockChainGetHeader] header:", header)
	avm.PushData(e, header)
	return true, nil
}

func (s *StateReader) BlockChainGetBlock(e *avm.ExecutionEngine) (bool, error) {
	data := avm.PopByteArray(e)
	var (
		block *ledger.Block
	)
	l := len(data)
	if l <= 5 {
		b := new(big.Int)
		height := uint32(b.SetBytes(data).Int64())
		if ledger.DefaultLedger != nil {
			hash, err := ledger.DefaultLedger.Store.GetBlockHash(height)
			if err != nil { return false, err }
			block, err = ledger.DefaultLedger.Store.GetBlock(hash)
			if err != nil { return false, err }
		}else {
			block = nil
		}
	}else if l == 32 {
		hash, err := common.Uint256ParseFromBytes(data)
		if err != nil { return false, err }
		if ledger.DefaultLedger != nil {
			block, err = ledger.DefaultLedger.Store.GetBlock(hash)
			if err != nil { return false, err }
		}else {
			block = nil
		}
	}else {
		return false, errors.New("The data length is error in function blockchaningetblock!")
	}
	avm.PushData(e, block)
	return true, nil
}

func (s *StateReader) BlockChainGetTransaction(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopByteArray(e)
	hash, err := common.Uint256ParseFromBytes(d)
	if err != nil {
		return false, err
	}
	tx, err := ledger.DefaultLedger.Store.GetTransaction(hash)
	if err != nil {
		return false, err
	}

	avm.PushData(e, tx)
	return true, nil
}

func (s *StateReader) BlockChainGetAccount(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopByteArray(e)
	hash, err := common.Uint160ParseFromBytes(d)
	if err != nil { return false, err }
	account, err := ledger.DefaultLedger.Store.GetAccount(hash)
	avm.PushData(e, account)
	return true, nil
}

func (s *StateReader) BlockChainGetAsset(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopByteArray(e)
	hash, err := common.Uint256ParseFromBytes(d)
	if err != nil { return false, err }
	asset, err := ledger.DefaultLedger.Store.GetAsset(hash)
	if err != nil { return false, err }
	avm.PushData(e, asset)
	return true, nil
}

func (s *StateReader) BlockChainGetValidators(e *avm.ExecutionEngine) (bool, error) {
	return true, nil
}

func (s *StateReader) HeaderGetHash(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get header error in function headergethash!", )
	}
   	h := d.(*ledger.Header).Blockdata.Hash()
	avm.PushData(e, h.ToArray())
	return true, nil
}

func (s *StateReader) HeaderGetVersion(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get header error in function headergetversion")
	}
	version := d.(*ledger.Header).Blockdata.Version
	avm.PushData(e, version)
	return true, nil
}

func (s *StateReader) HeaderGetPrevHash(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get header error in function HeaderGetPrevHash")
	}
	preHash := d.(*ledger.Header).Blockdata.PrevBlockHash
	avm.PushData(e, preHash.ToArray())
	return true, nil
}

func (s *StateReader) HeaderGetMerkleRoot(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get header error in function HeaderGetMerkleRoot")
	}
	root := d.(*ledger.Header).Blockdata.TransactionsRoot
	avm.PushData(e, root.ToArray())
	return true, nil
}

func (s *StateReader) HeaderGetTimestamp(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get header error in function HeaderGetTimestamp")
	}
	timeStamp := d.(*ledger.Header).Blockdata.Timestamp
	avm.PushData(e, timeStamp)
	return true, nil
}

func (s *StateReader) HeaderGetConsensusData(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get header error in function HeaderGetConsensusData")
	}
	consensusData := d.(*ledger.Header).Blockdata.ConsensusData
	avm.PushData(e, consensusData)
	return true, nil
}

func (s *StateReader) HeaderGetNextConsensus(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get header error in function HeaderGetNextConsensus")
	}
	nextBookKeeper := d.(*ledger.Header).Blockdata.NextBookKeeper
	avm.PushData(e, nextBookKeeper)
	return true, nil
}

func (s *StateReader) BlockGetTransactionCount(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get block error in function BlockGetTransactionCount")
	}
	transactions := d.(*ledger.Block).Transactions
	avm.PushData(e, len(transactions))
	return true, nil
}

func (s *StateReader) BlockGetTransactions(e *avm.ExecutionEngine) (bool, error) {
	fmt.Println("[BlockGetTransactions]")
	d := avm.PopInteropInterface(e)
	fmt.Println("[BlockGetTransactions] data", d)
	if d == nil {
		return false, fmt.Errorf("%v", "Get block data error in function BlockGetTransactions")
	}
	transactions := d.(*ledger.Block).Transactions
	transactionList := make([]types.StackItemInterface, 0)
	fmt.Println("================len transactions==============", len(transactions))
	for _, v := range transactions {
		transactionList = append(transactionList, types.NewInteropInterface(v))
	}
	fmt.Println("================transactionList==============", transactionList)
	avm.PushData(e, transactionList)
	return true, nil
}

func (s *StateReader) BlockGetTransaction(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get block data error in function BlockGetTransaction")
	}
	transactions := d.(*ledger.Block).Transactions
	index := avm.PopInt(e)
	avm.PushData(e, transactions[index])
	return true, nil
}

func (s *StateReader) TransactionGetHash(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get transaction error in function TransactionGetHash")
	}
	txHash := d.(*transaction.Transaction).Hash()
	avm.PushData(e, txHash.ToArray())
	return true, nil
}

func (s *StateReader) TransactionGetType(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get transaction error in function TransactionGetType")
	}
	txType := d.(*transaction.Transaction).TxType
	avm.PushData(e, int(txType))
	return true, nil
}

func (s *StateReader) TransactionGetAttributes(e *avm.ExecutionEngine) (bool, error) {
	return  true, nil
}

func (s *StateReader) TransactionGetInputs(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get transaction error in function TransactionGetInputs")
	}
	inputs := d.(*transaction.Transaction).UTXOInputs
	inputList := make([]types.StackItemInterface, 0)
	for _, v := range inputs {
		inputList = append(inputList, types.NewInteropInterface(v))
	}
	avm.PushData(e, inputList)
	return true, nil
}

func (s *StateReader) TransactionGetOutputs(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get transaction error in function TransactionGetOutputs")
	}
	outputs := d.(*transaction.Transaction).Outputs
	outputList := make([]types.StackItemInterface, 0)
	for _, v := range outputs {
		outputList = append(outputList, types.NewInteropInterface(v))
	}
	avm.PushData(e, outputList)
	return true, nil
}

func (s *StateReader) TransactionGetReferences(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get transaction error in function TransactionGetReferences")
	}
	references, err := d.(*transaction.Transaction).GetReference()
	referenceList := make([]types.StackItemInterface, 0)
	for _, v := range references {
		referenceList = append(referenceList, types.NewInteropInterface(v))
	}
	avm.PushData(e, referenceList)
	if err != nil {return false, err}
	return true, nil
}

func (s *StateReader) AttributeGetUsage(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get Attribute error in function AttributeGetUsage")
	}
	attribute := d.(*transaction.TxAttribute)
	avm.PushData(e, int(attribute.Usage))
	return true, nil
}

func (s *StateReader) AttributeGetData(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get Attribute error in function AttributeGetUsage")
	}
	attribute := d.(*transaction.TxAttribute)
	avm.PushData(e, attribute.Data)
	return true, nil
}

func (s *StateReader) InputGetHash(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get UTXOTxInput error in function InputGetHash")
	}
	input := d.(*transaction.UTXOTxInput)
	avm.PushData(e, input.ReferTxID.ToArray())
	return true, nil
}

func (s *StateReader) InputGetIndex(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get transaction error in function TransactionGetReferences")
	}
	input := d.(*transaction.UTXOTxInput)
	avm.PushData(e, input.ReferTxOutputIndex)
	return true, nil
}

func (s *StateReader) OutputGetAssetId(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get TxOutput error in function OutputGetAssetId")
	}
	output := d.(*transaction.TxOutput)
	avm.PushData(e, output.AssetID.ToArray())
	return true, nil
}

func (s *StateReader) OutputGetValue(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get TxOutput error in function OutputGetValue")
	}
	output := d.(*transaction.TxOutput)
	avm.PushData(e, output.Value.GetData())
	return true, nil
}

func (s *StateReader) OutputGetCodeHash(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get TxOutput error in function OutputGetCodeHash")
	}
	output := d.(*transaction.TxOutput)
	avm.PushData(e, output.ProgramHash.ToArray())
	return true, nil
}

func (s *StateReader) AccountGetCodeHash(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AccountState error in function AccountGetCodeHash")
	}
	accountState := d.(*states.AccountState).ProgramHash
	avm.PushData(e, accountState.ToArray())
	return true, nil
}

func (s *StateReader) AccountGetBalance(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AccountState error in function AccountGetCodeHash")
	}
	accountState := d.(*states.AccountState)
	assetIdByte := avm.PopByteArray(e)
	assetId, err := common.Uint256ParseFromBytes(assetIdByte)
	if err != nil {return false, err}
	balance := common.Fixed64(0)
	if v, ok := accountState.Balances[assetId]; ok {
		balance = v
	}
	avm.PushData(e, balance.GetData())
	return true, nil
}

func (s *StateReader) AssetGetAssetId(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AssetState error in function AssetGetAssetId")
	}
	assetState := d.(*states.AssetState)
	avm.PushData(e, assetState.AssetId.ToArray())
	return true, nil
}

func (s *StateReader) AssetGetAssetType(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AssetState error in function AssetGetAssetType")
	}
	assetState := d.(*states.AssetState)
	avm.PushData(e, int(assetState.AssetType))
	return true, nil
}

func (s *StateReader) AssetGetAmount(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AssetState error in function AssetGetAmount")
	}
	assetState := d.(*states.AssetState)
	avm.PushData(e, assetState.Amount.GetData())
	return true, nil
}

func (s *StateReader) AssetGetAvailable(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AssetState error in function AssetGetAvailable")
	}
	assetState := d.(*states.AssetState)
	avm.PushData(e, assetState.Available.GetData())
	return true, nil
}

func (s *StateReader) AssetGetPrecision(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AssetState error in function AssetGetPrecision")
	}
	assetState := d.(*states.AssetState)
	avm.PushData(e, int(assetState.Precision))
	return true, nil
}

func (s *StateReader) AssetGetOwner(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AssetState error in function AssetGetOwner")
	}
	assetState := d.(*states.AssetState)
	owner, err := assetState.Owner.EncodePoint(true)
	if err != nil {return false, err}
	avm.PushData(e, owner)
	return true, nil
}

func (s *StateReader) AssetGetAdmin(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AssetState error in function AssetGetAdmin")
	}
	assetState := d.(*states.AssetState)
	avm.PushData(e, assetState.Admin.ToArray())
	return true, nil
}

func (s *StateReader) AssetGetIssuer(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get AssetState error in function AssetGetIssuer")
	}
	assetState := d.(*states.AssetState)
	avm.PushData(e, assetState.Issuer.ToArray())
	return true, nil
}

func (s *StateReader) ContractGetCode(e *avm.ExecutionEngine) (bool, error) {
	d := avm.PopInteropInterface(e)
	if d == nil {
		return false, fmt.Errorf("%v", "Get ContractState error in function ContractGetCode")
	}
	assetState := d.(*states.ContractState)
	avm.PushData(e, assetState.Code.Code)
	return true, nil
}

func (s *StateReader) StorageGetContext(e *avm.ExecutionEngine) (bool, error) {
	codeHash, err := common.Uint160ParseFromBytes(e.CurrentContext().GetCodeHash())
	if err != nil {return false, err}
	avm.PushData(e, NewStorageContext(&codeHash))
	return true, nil
}
