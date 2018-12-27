package tentacle

import (
	"log"
	"bytes"
	"net/http"
	"github.com/gorilla/websocket"
	"jsonrpc"
)

var upgrader = websocket.Upgrader{}

func errorResponse(w http.ResponseWriter, r *http.Request, err error, status int, message string) {
	log.Printf("websocket error: %s %d", err.Error(), status)
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	context := Context()
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	actor := new(Actor).Init(conn)
	defer actor.Close()

	go actor.Start()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			// connection error, will close the connection
			//log.
			break
			//panic("close" + err.Error())
		}

		msg, err := jsonrpc.ParseMessage(data)
		if err != nil {
			errorResponse(w, r, err, 400, "Bad request")
			return
		}
		msg.FromConnId = actor.ConnId
		context.Router.RouteMessage(msg, actor.ConnId)
	}
}


func HandleHttp(w http.ResponseWriter, r *http.Request) {
	req := new(Requester).Init()
	defer req.Close()

	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(r.Body)
	if err != nil {
		errorResponse(w, r, err, 400, "Bad request")
		return
	}

	msg, err := jsonrpc.ParseMessage(buffer.Bytes())
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

