package httpwebsocket

import (
	. "DNA/common/config"
	"DNA/core/ledger"
	"DNA/events"
	"DNA/net/httprestful/common"
	Err "DNA/net/httprestful/error"
	"DNA/net/httpwebsocket/websocket"
	. "DNA/net/protocol"
	. "DNA/common"
)

var ws *websocket.WsServer
var pushBlockFlag bool = false

func StartServer(n Noder) {
	common.SetNode(n)
	ledger.DefaultLedger.Blockchain.BCEvents.Subscribe(events.EventBlockPersistCompleted, SendBlock2WSclient)
	go func() {
		ws = websocket.InitWsServer(common.CheckAccessToken)
		ws.Start()
	}()
}
func SendBlock2WSclient(v interface{}) {
	if Parameters.HttpWsPort != 0 && pushBlockFlag {
		go func() {
			PushBlock(v)
		}()
	}
}
func Stop() {
	if ws != nil {
		ws.Stop()
	}
}
func ReStartServer() {
	if ws == nil {
		ws = websocket.InitWsServer(common.CheckAccessToken)
		ws.Start()
		return
	}
	ws.Restart()
}
func GetWsPushBlockFlag() bool {
	return pushBlockFlag
}
func SetWsPushBlockFlag(b bool) {
	pushBlockFlag = b
}
func SetTxHashMap(txhash string, sessionid string) {
	if ws != nil {
		ws.SetTxHashMap(txhash, sessionid)
	}
}
func PushResult(txHash Uint256, errcode int64, action string, result interface{}) {
	if ws != nil {
		resp := common.ResponsePack(Err.SUCCESS)
		resp["Result"] = result
		resp["Error"] = errcode
		resp["Action"] = action
		resp["Desc"] = Err.ErrMap[resp["Error"].(int64)]
		ws.PushTxResult(ToHexString(txHash.ToArrayReverse()), resp)
	}
}
func PushBlock(v interface{}) {
	if ws != nil {
		resp := common.ResponsePack(Err.SUCCESS)
		if block, ok := v.(*ledger.Block); ok {
			resp["Result"] = common.GetBlockInfo(block)
			resp["Action"] = "sendrawblock"
			ws.PushResult(resp)
		}
	}
}
