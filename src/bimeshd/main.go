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
	fmt.Printf("version %s %s %d\n", cfg.Version, cfg.Bind.Host, cfg.Bind.Port)

	m.Print()

	StartServer()
}
