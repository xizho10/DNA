package net

import (
	. "DNA/common"
	"DNA/common/config"
	"DNA/core/ledger"
	"DNA/core/transaction"
	"DNA/crypto"
	"DNA/events"
	"DNA/net/node"
	"DNA/net/protocol"
)

type Neter interface {
	GetTxnPool(cleanPool bool) map[Uint256]*transaction.Transaction
	SynchronizeTxnPool()
	Xmit(interface{}) error
	GetEvent(eventName string) *events.Event
	GetBookKeepersAddrs() ([]*crypto.PubKey, uint64)
	CleanSubmittedTransactions(block *ledger.Block) error
	GetNeighborNoder() []protocol.Noder
	Tx(buf []byte)
}

func StartProtocol(pubKey *crypto.PubKey, nodeType int) (Neter, protocol.Noder) {
	seedNodes := config.Parameters.SeedList

	net := node.InitNode(pubKey, nodeType)
	for _, nodeAddr := range seedNodes {
		go net.Connect(nodeAddr)
	}
	return net, net
}
