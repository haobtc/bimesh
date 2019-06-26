package tentacle

import (
	"net/http"
	"github.com/gorilla/websocket"
	"jsonrpc"
)

var upgrader = websocket.Upgrader{}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	tentacle := Tentacle()
	conn, _ := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	actor := new(LocalConnT).Init(conn)
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
			jsonrpc.ErrorResponse(w, r, err, 400, "Bad request")
			return
		}
		msg.FromConnId = actor.ConnId
		tentacle.Router.RouteMessage(msg, actor.ConnId)
	}
}

func HandleHttp(w http.ResponseWriter, r *http.Request, msg jsonrpc.RPCMessage) (jsonrpc.RPCMessage, error){
	req := new(Requester).Init()
	defer req.Close()

	msg.FromConnId = req.ConnId
	result, err := req.RequestAndWait(msg)
	return result, err
}

