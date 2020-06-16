/**
 *
 * @author nghiatc
 * @since May 29, 2020
 */

package gserver

import (
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

const GSPrefix = ".gserver."

type GServer struct {
	Name string
	Addr string
	Conn net.Listener
	IsSSL bool
	Opts grpc.ServerOption
	Server *grpc.Server
}

func NewGServer(name string) *GServer {
	if len(name) == 0 {
		return nil
	}
	c := nconf.GetConfig()
	addr := c.GetString(name + GSPrefix + "address")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error server net listen tcp %v", err)
	}
	var s *grpc.Server
	var opts grpc.ServerOption
	isSSL := c.GetBool(name + GSPrefix + "is_ssl")
	if isSSL {
		certFile := c.GetString(name + GSPrefix + "cert_file") //"ssl/server.crt"
		keyFile := c.GetString(name + GSPrefix + "key_file")   //"ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Fatalf("Create creds ssl error %v\n", sslErr)
			return nil
		}
		opts = grpc.Creds(creds)
		s = grpc.NewServer(opts)
	} else {
		s = grpc.NewServer()
	}
	return &GServer{Name: name, Addr: addr, Conn: lis, IsSSL: isSSL, Opts: opts, Server: s}
}

func (gs *GServer) Start() error {
	fmt.Printf("GServer[%s] is running on: %v\n", gs.Name, gs.Addr)
	err := gs.Server.Serve(gs.Conn)
	if err != nil {
		log.Fatalf("GServer[%s] start error: %v\n", gs.Name, err)
	}
	return err
}

func (gs *GServer) Stop() {
	if gs != nil && gs.Server != nil {
		gs.Server.Stop()
	}
}
