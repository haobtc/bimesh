package main

import (
	"log"
	"fmt"
	"net/http"
	"bytes"
	"github.com/gorilla/websocket"
	"mesh"
	"datadir"
)

var upgrader = websocket.Upgrader{}

//var router = mesh.NewRouter()

// builtin services
//var serviceManager = new(mesh.ServiceManager)

func main() {
	//datadir.SetDataDir("hello")
	datadir.EnsureDataDir("")

	cfg := datadir.GetConfig()
	fmt.Printf("version %s %s %d\n", cfg.Version, cfg.Bind.Host, cfg.Bind.Port)

	router := datadir.GetRouter()
	for _, box := range router.Boxes {
		fmt.Printf("box %s %s %s \n", box.Boxid, box.Endpoint, box.Cert)
		for _, name := range box.ServiceNames {
			fmt.Printf(" - service %s\n", name)
		}
	}
	mesh.Context().Start()

	http.HandleFunc("/jsonrpc/ws", handleWebsocket)
	http.HandleFunc("/jsonrpc/http", handleHttp)

	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Bind.Host, cfg.Bind.Port), nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func errorResponse(w http.ResponseWriter, r *http.Request, err error, status int, message string) {
	log.Printf("HTTP error: %s %d", err.Error(), status)
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
	req := new(mesh.Requester).Init()
	defer req.Close()

	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		errorResponse(w, r, err, 400, "Bad request")
		return
	}

	msg, err := mesh.ParseMessage(buffer.Bytes())
	if err != nil {
		errorResponse(w, r, err, 400, "Bad request")
		return
	}

	msg.FromConnId = req.ConnId
	result, err := req.RequestAndWait(msg)
	if err != nil {
		errorResponse(w, r, err, 500, "server error")
		return
	}

	bytes, err := result.Raw.MarshalJSON()
	if err != nil {
		errorResponse(w, r, err, 500, "server error")
		return
	}
	w.Write(bytes)
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	context := mesh.Context()
	c, _ := upgrader.Upgrade(w, r, nil)
	defer c.Close()

	actor := new(mesh.Actor).Init(c)
	defer actor.Close()

	go actor.Start()

	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			// connection error, will close the connection
			//log.
			break
			//panic("close" + err.Error())
		}

		msg, err := mesh.ParseMessage(data)
		if err != nil {
			errorResponse(w, r, err, 400, "Bad request")
			return
		}
		msg.FromConnId = actor.ConnId
		context.Router.RouteMessage(msg, actor.ConnId)
	}
}
