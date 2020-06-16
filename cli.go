/**
 *
 * @author nghiatc
 * @since Jun 16, 2020
 */
package main

import (
	"context"
	"fmt"
	"github.com/congnghia0609/ntc-gconf/nconf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"ntc-ggrpc/example/ngrpc"
	"ntc-ggrpc/gclient"
	"path/filepath"
	"runtime"
	"time"
)

func InitNConf2() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

func main() {
	// Init NConf
	InitNConf2()

	// Calculator Client
	name := "ngrpc"
	cc := NewCalClient(name)

	cc.CallSum()

	cc.CallSumWithDeadline(1 * time.Second)
	cc.CallSumWithDeadline(5 * time.Second)

	cc.CallPND()

	cc.CallAverage()

	cc.CallFindMax()

	cc.CallSquareRoot(9)
}

type CalClient struct {
	Name string
	GCli *gclient.GClient
	Cli  ngrpc.CalculatorServiceClient
}

func NewCalClient(name string) *CalClient {
	gcli := gclient.NewGClient(name)
	cli := ngrpc.NewCalculatorServiceClient(gcli.Conn)
	return &CalClient{Name: name, GCli: gcli, Cli: cli}
}

func (cc *CalClient) Close() {
	if cc != nil && cc.GCli != nil {
		cc.GCli.Close()
	}
}

func (cc *CalClient) CallSum() {
	log.Println("Call sum api")
	resp, err := cc.Cli.Sum(context.Background(), &ngrpc.SumRequest{
		Num1: 7,
		Num2: 6,
	})
	if err != nil {
		log.Fatalf("call sum api error %v\n", err)
	}
	log.Printf("sum api response %v\n", resp.GetResult())
}

func (cc *CalClient) CallSumWithDeadline(timeout time.Duration) {
	log.Println("Call sum with deadline api")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := cc.Cli.SumWithDeadline(ctx, &ngrpc.SumRequest{
		Num1: 7,
		Num2: 6,
	})
	if err != nil {
		if statusErr, ok := status.FromError(err); ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("call sum with deadline DeadlineExceeded")
			} else {
				log.Fatalf("call sum with deadline api error %v\n", err)
			}
		} else {
			log.Fatalf("call sum with deadline api unknown error %v\n", err)
		}
		return
	}
	log.Printf("sum with deadline api response %v\n", resp.GetResult())
}

func (cc *CalClient) CallPND() {
	log.Println("call pnd api")
	stream, err := cc.Cli.PrimeNumberDecomposition(context.Background(), &ngrpc.PNDRequest{
		Number: 120,
	})
	if err != nil {
		log.Fatalf("callPND err %v\n", err)
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("server finish streaming")
			return
		}
		if err != nil {
			log.Fatalf("callPND error %v\n", err)
		}

		log.Printf("prime number %v\n", resp.GetResult())
	}
}

func (cc *CalClient) CallAverage() {
	log.Println("call average api")
	stream, err := cc.Cli.Average(context.Background())
	if err != nil {
		log.Fatalf("call average error %v\n", err)
	}

	listReq := []ngrpc.AverageRequest{
		ngrpc.AverageRequest{Num: 5},
		ngrpc.AverageRequest{Num: 10},
		ngrpc.AverageRequest{Num: 15},
		ngrpc.AverageRequest{Num: 20},
		ngrpc.AverageRequest{Num: 25},
	}
	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("Send average request error %v\n", err)
		}
		time.Sleep(1 * time.Second)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("receive average response error %v", err)
	}
	log.Printf("Average response %v\n", resp)
}

func (cc *CalClient) CallFindMax() {
	log.Println("call find max")
	stream, err := cc.Cli.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call find max error %v\n", err)
	}
	waitc := make(chan struct{})
	go func() {
		// Gui nhieu request
		listReq := []ngrpc.FindMaxRequest{
			ngrpc.FindMaxRequest{Num: 5},
			ngrpc.FindMaxRequest{Num: 10},
			ngrpc.FindMaxRequest{Num: 1},
			ngrpc.FindMaxRequest{Num: 6},
			ngrpc.FindMaxRequest{Num: 9},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("send find max request error %v\n", err)
				break
			}
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		defer close(waitc)
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("End find max api.")
				break
			}
			if err != nil {
				log.Fatalf("revc find max error %v", err)
				break
			}
			log.Printf("max: %v\n", resp.GetMax())
		}
	}()

	<-waitc
}

func (cc *CalClient) CallSquareRoot(num int32) {
	log.Println("call square root api")
	resp, err := cc.Cli.Square(context.Background(), &ngrpc.SquareRequest{
		Num: num,
	})
	if err != nil {
		log.Printf("call square root api error %v\n", err)
		if errStatus, ok := status.FromError(err); ok {
			log.Printf("error code: %v\n", errStatus.Code())
			log.Printf("error msg: %v\n", errStatus.Message())
			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("InvalidArgument num %v\n", num)
				return
			}
		}
	}
	log.Printf("square root response %v\n", resp.GetSquareRoot())
}
