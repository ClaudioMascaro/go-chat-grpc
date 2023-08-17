package main

import (
	"log"
	"net"
	"os"

	"go-cli-chat-client/chatservice"

	"google.golang.org/grpc"
)

func main() {
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "50051"
	}

	listen, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ : " + Port)

	grpcServer := grpc.NewServer()

	cs := chatservice.ChatService{}
	chatservice.RegisterChatServer(grpcServer, &cs)

	err = grpcServer.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to start gRPC Server :: %v", err)
	}
}
