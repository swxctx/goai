package main

import (
	"github.com/swxctx/goai/example/api"
	td "github.com/swxctx/malatd"
)

func main() {
	// Gen Time: 2024-04-15 18:08:50
	srv := td.NewServer(cfg.SrvConfig)
	api.Route(srv, "/example")
	srv.Run()
}
