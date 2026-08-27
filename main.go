package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"raft-consensus/handler"
	pb "raft-consensus/proto"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	id := flag.Int("id", 0, "Node ID")
	port := flag.Int("port", 0, "Listening port")
	peersRaw := flag.String("peers", "", "Comma-separated peer addresses")

	flag.Parse()

	if *id == 0 {
		log.Fatal("No id specified")
	}
	if *port == 0 {
		log.Fatal("No port specified")
	}
	if *peersRaw == "" {
		log.Fatal("No peers specified")
	}

	peers := strings.Split(*peersRaw, ",")
	_ = peers

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serv := grpc.NewServer()
	handler := handler.RaftServer{}
	pb.RegisterRaftServer(serv, &handler)

	log.Printf("Node %d is serving on port %d", *id, *port)
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
