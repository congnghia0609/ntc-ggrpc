# ntc-ggrpc
ntc-ggrpc is an example golang gRPC  


## Quick start
```bash
// Install dependencies
//make deps
go mod download

// update go.mod file
go mod tidy

// Gen Source Code gRPC
make gen

// Gen resource SSL
make ssl

// Run Server
make server

// Run client
make client
```


## gRPC Server
Example CalculatorService Server:  
```go
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
```

## gRPC Client
Example CalculatorService Client  
```go
func main() {
	// Init NConf
	InitNConf()

	// Calculator Client
	name := "ngrpc"
	cc := NewCalClient(name)
	defer cc.Close()

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
```


## HAProxy Config Load Balancer for gRPC Server
```bash
frontend grpc_fe
	bind *:3330 ssl crt /ssl/haproxy.pem alpn h2
	option http-use-htx
    option logasap
	default_backend grpc_be

backend grpc_be
	balance roundrobin
	server grpc-go 127.0.0.1:3333 check ssl alpn h2 ca-file /ssl/server.crt
	server grpc-java 127.0.0.1:3334 check ssl alpn h2 ca-file /ssl/server.crt
```

## NginX Config Load Balancer for gRPC Server, mode Non-SSL
Reference document at [NginX blog](https://www.nginx.com/blog/nginx-1-13-10-grpc/) and [Module ngx_http_grpc_module](http://nginx.org/en/docs/http/ngx_http_grpc_module.html)  
```bash
http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  1800;

    map $http_upgrade $connection_upgrade {
         default upgrade;
         ''      close;
    }

    upstream grpcservers {
        server 127.0.0.1:3333;
        server 127.0.0.1:3334;
    }

    server {
        listen 3330 http2;

        # Can use location detail with a path of gRPC Service. Ex: /helloworld.Greeter
        location /ngrpc.CalculatorService {
            # The 'grpc://' ==> Non-SSL gRPC
            # The 'grpcs://' ==> SSL gRPC
            grpc_pass grpc://grpcservers;
        }
    }
}
```

## NginX Config Load Balancer for gRPC Server, mode SSL
```bash
http {
    include       mime.types;
    default_type  application/octet-stream;
    sendfile        on;
    keepalive_timeout  1800;

    map $http_upgrade $connection_upgrade {
         default upgrade;
         ''      close;
    }

    upstream grpcservers {
        server 127.0.0.1:3333;
        server 127.0.0.1:3334;
    }

    server {
        listen 3330 ssl http2;

        ## Create a certificate that points to the hostname.
        ## $ openssl req -nodes -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -subj '/CN=nginx'
        ssl_certificate /ssl/server.crt;
        ssl_certificate_key /ssl/server.pem;

        # Can use location detail with a path of gRPC Service. Ex: /helloworld.Greeter
        location /ngrpc.CalculatorService {
            # The 'grpc://' ==> Non-SSL gRPC
            # The 'grpcs://' ==> SSL gRPC
            grpc_pass grpcs://grpcservers;
        }
    }
}
```

## License
This code is under the [Apache License v2](https://www.apache.org/licenses/LICENSE-2.0).  
