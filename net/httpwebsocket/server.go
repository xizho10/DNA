package httpwebsocket

import (
	. "DNA/common"
	. "DNA/common/config"
	"DNA/common/log"
	"DNA/core/ledger"
	"DNA/events"
	. "DNA/net/httprestful/common"
	Err "DNA/net/httprestful/error"
	"DNA/net/httpwebsocket/websocket"
	. "DNA/net/protocol"
)

const OAUTH_SUCCESS_CODE = "r0000"

var ws *websocket.WsServer
var pushBlockFlag bool = false

func StartServer(n Noder) {
	SetNode(n)
	ledger.DefaultLedger.Blockchain.BCEvents.Subscribe(events.EventBlockPersistCompleted, SendBlock2WSclient)
	go func() {
		ws = websocket.InitWsServer(checkAccessToken)
		ws.Start()
	}()
}
func SetWsPushBlockFlag(b bool) {
	pushBlockFlag = b
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
		ws = websocket.InitWsServer(checkAccessToken)
		ws.Start()
		return
	}
	ws.Restart()
}

func SetTxHashMap(txhash string, sessionid string) {
	if ws != nil {
		ws.SetTxHashMap(txhash, sessionid)
	}
}
func PushSmartCodeInvokeResult(txHash Uint256, errcode int64, result interface{}) {
	if ws != nil {
		resp := ResponsePack(Err.SUCCESS)
		var Result = make(map[string]interface{})
		txHashStr := ToHexString(txHash.ToArray())
		Result["TxHash"] = txHashStr
		Result["ExecResult"] = result

		resp["Result"] = Result
		resp["Action"] = "sendsmartcodeinvoke"
		resp["Error"] = errcode
		resp["Desc"] = Err.ErrMap[resp["Error"].(int64)]
		ws.PushTxResult(txHashStr, resp)
	}
}
func PushBlock(v interface{}) {
	if ws != nil {
		resp := ResponsePack(Err.SUCCESS)
		if block, ok := v.(*ledger.Block); ok {
			resp["Result"] = GetBlockInfo(block)
			resp["Action"] = "pushblock"
			ws.PushResult(resp)
		}
	}
}
func checkAccessToken(auth_type, access_token string) (cakey string, errCode int64, result interface{}) {

	if len(Parameters.OauthServerAddr) == 0 {
		return "", Err.SUCCESS, ""
	}
	req := make(map[string]interface{})
	req["token"] = access_token
	req["auth_type"] = auth_type
	repMsg, err := OauthRequest("GET", req, Parameters.OauthServerAddr)
	if err != nil {
		log.Error("Oauth timeout:", err)
		return "", Err.OAUTH_TIMEOUT, repMsg
	}
	if repMsg["code"] == OAUTH_SUCCESS_CODE {
		msg, ok := repMsg["msg"].(map[string]interface{})
		if !ok {
			return "", Err.INVALID_TOKEN, repMsg
		}
		if CAkey, ok := msg["cakey"].(string); ok {
			return CAkey, Err.SUCCESS, repMsg
		}
	}
	return "", Err.INVALID_TOKEN, repMsg
}
