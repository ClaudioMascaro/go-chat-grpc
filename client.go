package main

import (
	"bufio"
	"context"
	"fmt"
	"go-cli-chat-client/chatservice"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	serverAddress := os.Getenv("SERVER_ID")
	serverAddress = strings.Trim(serverAddress, "\r\n")

	log.Println("Connecting : " + serverAddress)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Faile to conncet to gRPC server :: %v", err)
	}
	defer conn.Close()

	client := chatservice.NewChatClient(conn)

	stream, err := client.StreamMessages(context.Background())
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	ch := clientHandle{stream: stream}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	bl := make(chan bool)
	<-bl

}

type clientHandle struct {
	stream     chatservice.Chat_StreamMessagesClient
	clientName string
}

func (ch *clientHandle) clientConfig() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf(" Failed to read from console :: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")

}

func (ch *clientHandle) sendMessage() {
	for {

		reader := bufio.NewReader(os.Stdin)
		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf(" Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chatservice.Message{
			User: ch.clientName,
			Text: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)

		if err != nil {
			log.Printf("Error while sending message to server :: %v", err)
		}

	}

}

func (ch *clientHandle) receiveMessage() {
	for {
		msg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving message from server :: %v", err)
		}

		fmt.Printf("%s : %s \n", msg.User, msg.Text)

	}
}
