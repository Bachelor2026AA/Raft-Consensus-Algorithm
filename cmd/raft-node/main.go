package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"raft-consensus"
	pb "raft-consensus/proto"
	"raft-consensus/repository"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	id, port, peers := parseFlags()

	clients := peerClients(peers)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logEntryRepo := &repository.LogEntryRepositoryImpl{}
	keyValueRepo := &repository.KeyValueRepositoryImpl{}

	serv := grpc.NewServer()
	node := raft.NewNode(id, logEntryRepo, keyValueRepo, clients)
	pb.RegisterRaftServer(serv, node)

	log.Printf("Node %d is serving on port %d", id, port)
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func parseFlags() (int, int, string) {
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
	return *id, *port, *peersRaw
}

func peerClients(peers string) []pb.RaftClient {
	var peerClients []pb.RaftClient
	for _, peer := range strings.Split(peers, ",") {
		conn, err := grpc.NewClient(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()
		peerClients = append(peerClients, pb.NewRaftClient(conn))
	}
	return peerClients
}
