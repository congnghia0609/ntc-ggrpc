# ntc-ggrpc
ntc-ggrpc is a example golang gRPC

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
