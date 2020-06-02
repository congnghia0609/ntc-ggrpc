/**
 *
 * @author nghiatc
 * @since May 29, 2020
 */

package main

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"ntc-ggrpc/ngrpc"
	"time"
)

func main() {
	target := "localhost:3330" // grpc-haproxy
	//target := "localhost:3333" // grpc-go

	//====== Begin Mode SSL ======//
	certFile := "ssl/client.crt"
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	if sslErr != nil {
		log.Fatalf("Create client creds ssl error %v\n", sslErr)
		return
	}
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(creds))
	//====== End Mode SSL ======//

	//conn, err := grpc.Dial(target, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		log.Fatalf("Error dial connection %v", err)
	}

	client := ngrpc.NewCalculatorServiceClient(conn)

	//callSum(client)

	//callSumWithDeadline(client, 1*time.Second)
	//callSumWithDeadline(client, 5*time.Second)

	//callPND(client)

	//callAverage(client)

	callFindMax(client)

	callSquareRoot(client, 9)
}

func callSum(c ngrpc.CalculatorServiceClient) {
	log.Println("Call sum api")
	resp, err := c.Sum(context.Background(), &ngrpc.SumRequest{
		Num1: 7,
		Num2: 6,
	})
	if err != nil {
		log.Fatalf("call sum api error %v\n", err)
	}
	log.Printf("sum api response %v\n", resp.GetResult())
}

func callSumWithDeadline(c ngrpc.CalculatorServiceClient, timeout time.Duration) {
	log.Println("Call sum with deadline api")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := c.SumWithDeadline(ctx, &ngrpc.SumRequest{
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

func callPND(c ngrpc.CalculatorServiceClient) {
	log.Println("call pnd api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &ngrpc.PNDRequest{
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

func callAverage(c ngrpc.CalculatorServiceClient) {
	log.Println("call average api")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("call average error %v\n", err)
	}

	listReq := []ngrpc.AverageRequest {
		ngrpc.AverageRequest{ Num: 5, },
		ngrpc.AverageRequest{ Num: 10, },
		ngrpc.AverageRequest{ Num: 15, },
		ngrpc.AverageRequest{ Num: 20, },
		ngrpc.AverageRequest{ Num: 25, },
	}
	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("Send average request error %v\n", err)
		}
		time.Sleep(1*time.Second)
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("receive average response error %v", err)
	}
	log.Printf("Average response %v\n", resp)
}

func callFindMax(c ngrpc.CalculatorServiceClient) {
	log.Println("call find max")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("call find max error %v\n", err)
	}
	waitc := make(chan struct{})
	go func() {
		// Gui nhieu request
		listReq := []ngrpc.FindMaxRequest{
			ngrpc.FindMaxRequest{Num: 5,},
			ngrpc.FindMaxRequest{Num: 10,},
			ngrpc.FindMaxRequest{Num: 1,},
			ngrpc.FindMaxRequest{Num: 6,},
			ngrpc.FindMaxRequest{Num: 9,},
		}
		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalf("send find max request error %v\n", err)
				break
			}
			time.Sleep(1*time.Second)
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

func callSquareRoot(c ngrpc.CalculatorServiceClient, num int32) {
	log.Println("call square root api")
	resp, err := c.Square(context.Background(), &ngrpc.SquareRequest{
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
