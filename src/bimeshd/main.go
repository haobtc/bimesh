package main

import (
	"fmt"
	"datadir"
	"mesh"
	"bbox"
)

func main() {
	//datadir.SetDataDir("hello")
	datadir.EnsureDataDir("")

	m := mesh.GetMesh()
	cfg := datadir.GetConfig()
	fmt.Printf("version %s %s %d\n", cfg.Version, cfg.Bind.Host, cfg.Bind.Port)

	for _, endp := range cfg.Static.Endpoints {
		fmt.Printf("endpoint %s %s %s\n", endp.Url, endp.ServiceType, endp.Cert)
		for _, name := range endp.ServiceNames {
			fmt.Printf(" - service %s\n", name)
		}
	}

	m.Print()

	StartServer()
}
