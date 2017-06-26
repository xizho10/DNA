package main

import (
	"DNA/net/httprestful/restful"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	//	"time"
	"sync"
)

var mux sync.Mutex
var i int = 0

// 定义http请求的处理方法
func handlerHello(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}
	request := make(map[string]interface{})
	err = json.Unmarshal(body, &request)
	if err != nil {

	}
	mux.Lock()
	i++
	fmt.Println(i)
	mux.Unlock()
	//fmt.Println(r.URL.Path,request)
	resp := map[string]interface{}{
		"Action":  "zx",
		"Result":  "restful-listen-dsggshshhhshitirurteetwt",
		"Error":   0,
		"Desc":    "",
		"Version": "1.0.0",
	}
	ret, _ := json.Marshal(resp)
	//time.Sleep(time.Second / 500)
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(ret))
}

func main() {

	dnaFlag := false
	if dnaFlag {
		rest := restful.InitRestServer(checkAccessToken)
		rest.Start()
		return
	}

	router := restful.NewRouter()
	//router.Get("/", handlerHello)
	router.Post("/", handlerHello)
	listener, err := net.Listen("tcp", ":60334")
	if err != nil {
	}

	server := &http.Server{Handler: router}
	fmt.Println("Server on 60334 start.")
	err = server.Serve(listener)
	if err != nil {
	}

	os.Exit(0)
}
func checkAccessToken(auth_type, access_token string) (cakey string, errCode int64, result interface{}) {
	return "", 0, ""
}
