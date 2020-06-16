/**
 *
 * @author nghiatc
 * @since May 29, 2020
 */

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
	"ntc-ggrpc/example/ngrpc"
	"time"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *ngrpc.SumRequest) (*ngrpc.SumResponse, error) {
	log.Println("sum called...")
	resp := &ngrpc.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}
	return resp, nil
}

func (s *server) SumWithDeadline(ctx context.Context, req *ngrpc.SumRequest) (*ngrpc.SumResponse, error) {
	log.Println("sum with deadline called...")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("context.Canceled...")
			return nil, status.Error(codes.Canceled, "client canceled req")
		}
		time.Sleep(1 * time.Second)
	}

	resp := &ngrpc.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func (s *server) PrimeNumberDecomposition(req *ngrpc.PNDRequest, stream ngrpc.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Println("PrimeNumberDecomposition called...")
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			// send to client
			stream.Send(&ngrpc.PNDResponse{
				Result: k,
			})
			log.Printf("server response k = %v\n", k)
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
			//log.Printf("k increase to %v\n", k)
		}
	}
	return nil
}

func (s *server) Average(stream ngrpc.CalculatorService_AverageServer) error {
	log.Println("Average called...")
	var total = float32(0)
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Tinh trung binh va return ve client
			resp := &ngrpc.AverageResponse{
				Result: total / float32(count),
			}
			return stream.SendAndClose(resp)
		}
		if err != nil {
			log.Fatalf("Average error %v\n", err)
			return err
		}
		log.Printf("receive req %v\n", req)
		total += req.GetNum()
		count++
	}
}

func (s *server) FindMax(stream ngrpc.CalculatorService_FindMaxServer) error {
	log.Println("FindMax called...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("FindMax EOF...")
			return nil
		}
		if err != nil {
			log.Fatalf("FindMax recv error %v\n", err)
			return err
		}
		num := req.GetNum()
		log.Printf("recv num %v\n", num)
		if num > max {
			max = num
		}
		err = stream.Send(&ngrpc.FindMaxResponse{
			Max: max,
		})
		if err != nil {
			log.Fatalf("Send mac error %v\n", err)
			return err
		}
	}
}

func (s *server) Square(ctx context.Context, req *ngrpc.SquareRequest) (*ngrpc.SquareResponse, error) {
	log.Println("Square called...")
	num := req.GetNum()
	if num <= 0 {
		log.Printf("InvalidArgument num < 0, num=%v\n", num)
		return nil, status.Errorf(codes.InvalidArgument, "Expect num > 0, req num was %v\n", num)
	}
	return &ngrpc.SquareResponse{
		SquareRoot: math.Sqrt(float64(num)),
	}, nil
}

func main() {
	addr := "0.0.0.0:3333"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error server net listen tcp %v", err)
	}

	//====== Begin Mode SSL ======//
	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Create creds ssl error %v\n", sslErr)
		return
	}
	opts := grpc.Creds(creds)
	s := grpc.NewServer(opts)
	//====== End Mode SSL ======//

	//s := grpc.NewServer()
	ngrpc.RegisterCalculatorServiceServer(s, &server{})

	fmt.Printf("Server is running on: %v\n", addr)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Server start error %v", err)
	}
}
