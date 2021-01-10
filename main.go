/**
 *
 * @author nghiatc
 * @since Jun 16, 2020
 */
package main

import (
	"fmt"
	"ntc-ggrpc/example/ghandler"
	"ntc-ggrpc/example/ngrpc"
	"ntc-ggrpc/gserver"
	"path/filepath"
	"runtime"

	"github.com/congnghia0609/ntc-gconf/nconf"
)

func InitNConf() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

func StartCalServer() {
	name := "ngrpc"
	gs := gserver.NewGServer(name)
	ngrpc.RegisterCalculatorServiceServer(gs.Server, &ghandler.CalculatorHandler{})
	gs.Start()
}

func main() {
	// Init NConf
	InitNConf()

	// Start CalServer
	StartCalServer()
}
