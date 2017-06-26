package main

import (
	"net/http"
	"sync"
	//"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var mainMux ServeMux

//multiplexer that keeps track of every function to be called on specific rpc call
type ServeMux struct {
	sync.RWMutex
	m               map[string]func([]interface{}) map[string]interface{}
	defaultFunction func(http.ResponseWriter, *http.Request)
}

//a function to register functions to be called for specific rpc calls
func HandleFunc(pattern string, handler func([]interface{}) map[string]interface{}) {
	mainMux.Lock()
	defer mainMux.Unlock()
	mainMux.m[pattern] = handler
}

var i int = 0

func Handle(w http.ResponseWriter, r *http.Request) {
	//mainMux.RLock()
	//defer mainMux.RUnlock()

	if r.Body == nil {
		if mainMux.defaultFunction != nil {

			mainMux.defaultFunction(w, r)
		}
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {

	}
	request := make(map[string]interface{})
	err = json.Unmarshal(body, &request)
	if err != nil {

	}
	i++
	fmt.Println(i)
	//fmt.Println(r.URL.Path,request)
	var resp map[string]interface{}
	if request["method"] != nil {
		//function, _ := mainMux.m[request["method"].(string)]
		//resp = function([]interface{}{})
	}
	resp = handlerHello([]interface{}{})

	ret, _ := json.Marshal(resp)
	//time.Sleep(time.Second / 500)
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Connection", "close")
	w.Write([]byte(ret))
}
func handlerHello(params []interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"Action":  "zx-rpc",
		"Result":  "rpc-dsggshshhhshitirurteetwt",
		"Error":   0,
		"Desc":    "",
		"Version": "1.0.0",
	}
	return resp

}
func main() {
	mainMux.m = make(map[string]func([]interface{}) map[string]interface{})
	http.HandleFunc("/", Handle)
	HandleFunc("getblock", handlerHello)
	fmt.Println("Server on 50336 start.")
	http.ListenAndServe(":50336", nil)
}
