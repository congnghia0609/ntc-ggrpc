/**
 *
 * @author nghiatc
 * @since May 29, 2020
 */

package gclient

import (
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

const GCPrefix = ".gclient."

type GClient struct {
	Name string
	Target string
	IsSSL bool
	Conn *grpc.ClientConn
}

func NewGClient(name string) *GClient {
	if len(name) == 0 {
		return nil
	}
	c := nconf.GetConfig()
	target := c.GetString(name + GCPrefix + "target") //"localhost:3333"
	isSSL := c.GetBool(name + GCPrefix + "is_ssl")
	fmt.Printf("NewGClient isSSL: %v\n", isSSL)
	var conn *grpc.ClientConn
	var err error
	if isSSL {
		certFile := c.GetString(name + GCPrefix + "cert_file") //"ssl/client.crt"
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Create client creds ssl error %v\n", sslErr)
			return nil
		}
		conn, err = grpc.Dial(target, grpc.WithTransportCredentials(creds))
		//====== End Mode SSL ======//
	} else {
		conn, err = grpc.Dial(target, grpc.WithInsecure())
	}
	if err != nil {
		log.Fatalf("Error dial connection %v", err)
	}
	return &GClient{Name: name, Target: target, IsSSL: isSSL, Conn: conn}
}

func (gc *GClient) Close() {
	if gc != nil && gc.Conn != nil {
		gc.Conn.Close()
	}
}
