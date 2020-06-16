package main

import (
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"ntc-ggrpc/ghandler"
	"ntc-ggrpc/ngrpc"
	"ntc-ggrpc/gserver"
	"path/filepath"
	"runtime"
)

func InitNConf() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

func main() {
	// Init NConf
	InitNConf()

	// Start GServer
	name := "ngrpc"
	gs := gserver.NewGServer(name)
	ngrpc.RegisterCalculatorServiceServer(gs.Server, &ghandler.CalculatorHandler{})
	gs.Start()
}