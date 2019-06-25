package main

import (
	"log"
	"fmt"
	"bytes"
	"errors"
	"net/http"
	"tentacle"
	"datadir"
	"mesh"
	"jsonrpc"
)

func StartServer() {
	cfg := datadir.GetConfig()
	
	tentacle.Context().Start()

	http.HandleFunc("/jsonrpc/ws", HandleWebsocket)
	http.HandleFunc("/jsonrpc/http", HandleHttp)

	http.HandleFunc("/", HandleHome)
	log.Fatal(
		http.ListenAndServe(
		fmt.Sprintf("%s:%d", cfg.Bind.Host, cfg.Bind.Port), nil))
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	// currently only tentacle module support websocket connection
	tentacle.HandleWebsocket(w, r)
}

func HandleHttp(w http.ResponseWriter, r*http.Request) {
	// currently only tentacle module support http connection
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		jsonrpc.ErrorResponse(w, r, err, 400, "Bad request")
		return
	}

	msg, err := jsonrpc.ParseMessage(buffer.Bytes())
	if err != nil {
		jsonrpc.ErrorResponse(w, r, err, 400, "Bad request")
		return
	}

	//result, err := tentacle.HandleHttp(w, r, msg)
	//endpoint = mesh.GetMesh().GetEndpoint()
	if msg.ServiceName == "" {
		jsonrpc.ErrorResponse(w, r, errors.New("bad or nil service name"), 400, "Bad request")
		return
	}
	endpoint := mesh.GetMesh().GetEndpoint(msg.ServiceName)
	if endpoint == nil {
		jsonrpc.ErrorResponse(w, r, errors.New("service not found"), 404, "Not Found")
		return
	}

	result, err := (*endpoint).Request(msg)
	if err != nil {
		jsonrpc.ErrorResponse(w, r, err, 500, "Server error")
	}
	data, err := result.Raw.MarshalJSON()
	if err != nil {
		jsonrpc.ErrorResponse(w, r, err, 500, "Server error")
	}
	w.Write(data)
}
