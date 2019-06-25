package main

import (
	"fmt"
	"datadir"
	"mesh"
)

func main() {
	//datadir.SetDataDir("hello")
	datadir.EnsureDataDir("")
	m := mesh.GetMesh()
	cfg := datadir.GetConfig()
	fmt.Printf("version %s %s\n", cfg.Version, cfg.Server.Bind)

	m.Print()

	StartServer()
}
