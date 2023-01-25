package main

import (
	models "SheetUsingProtoBuffer/model" //import root folder of grpc
	"google.golang.org/grpc"
	"log"
	"net"
)

type api struct {
	models.UnimplementedAPIServer
}

func main() {
	//for request at server, it listens the TCP protocol at :443 port
	lis, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatalf("failed to listen :443") //here, we don't use http coz in grpc http2 protocol is used
		//also, in GRPC, we never return anything
	}
	//now, create new grpcserver to call new grpc protocol
	var grpcServer = grpc.NewServer()

}
