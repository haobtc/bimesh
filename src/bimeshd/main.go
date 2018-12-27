package main

import (
	"log"
	"fmt"
	"net/http"
	"tentacle"
	"datadir"
)

func main() {
	//datadir.SetDataDir("hello")
	datadir.EnsureDataDir("")

	cfg := datadir.GetConfig()
	fmt.Printf("version %s %s %d\n", cfg.Version, cfg.Bind.Host, cfg.Bind.Port)

	for _, box := range cfg.StaticBoxes {
		fmt.Printf("box %s %s %s \n", box.Boxid, box.Endpoint, box.Cert)
		for _, name := range box.ServiceNames {
			fmt.Printf(" - service %s\n", name)
		}
	}

	tentacle.Context().Start()

	http.HandleFunc("/jsonrpc/ws", HandleWebsocket)
	http.HandleFunc("/jsonrpc/http", HandleHttp)

	http.HandleFunc("/", home)
	log.Fatal(
		http.ListenAndServe(
		fmt.Sprintf("%s:%d", cfg.Bind.Host, cfg.Bind.Port), nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	// currently only tentacle module support websocket connection
	tentacle.HandleWebsocket(w, r)
}

func HandleHttp(w http.ResponseWriter, r*http.Request) {
	// currently only tentacle module support http connection
	tentacle.HandleHttp(w, r)
}
