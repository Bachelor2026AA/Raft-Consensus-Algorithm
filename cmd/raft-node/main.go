package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"raft-consensus"
	"raft-consensus/repository"

	pb "raft-consensus/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	id, grpcPort, httpPort, peers := parseFlags()

	clients := map[int]pb.RaftClient{}
	for i, peer := range strings.Split(peers, ",") {
		conn, err := grpc.NewClient(peer, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()
		i++
		clients[i] = pb.NewRaftClient(conn)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logEntryRepo := &repository.LogEntryRepositoryImpl{}
	keyValueRepo := &repository.KeyValueRepositoryImpl{}

	serv := grpc.NewServer()
	node := raft.NewNode(id, logEntryRepo, keyValueRepo, clients)
	pb.RegisterRaftServer(serv, node)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /keys/{id}", node.Get)
	mux.HandleFunc("POST /keys", node.Create)
	mux.HandleFunc("PATCH /keys/{id}", node.Update)
	mux.HandleFunc("DELETE /keys/{id}", node.Delete)
	go func() {
		log.Printf("Node %d is serving http on port %d", id, httpPort)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), mux); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Printf("Node %d is serving grpc on port %d", id, grpcPort)
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func parseFlags() (int, int, int, string) {
	id := flag.Int("id", 0, "Node ID")
	grpcPort := flag.Int("grpc-port", 0, "GRPC listening port")
	httpPort := flag.Int("http-port", 0, "HTTP listening port")
	peersRaw := flag.String("peers", "", "Comma-separated peer addresses")

	flag.Parse()

	if *id == 0 {
		log.Fatal("No id specified")
	}
	if *grpcPort == 0 {
		log.Fatal("No GRPC port specified")
	}
	if *httpPort == 0 {
		log.Fatal("No HTTP port specified")
	}
	if *peersRaw == "" {
		log.Fatal("No peers specified")
	}
	return *id, *grpcPort, *httpPort, *peersRaw
}
