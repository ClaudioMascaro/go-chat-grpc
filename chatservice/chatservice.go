package chatservice

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode string
	ClientUniqueCode  string
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageHandleObject = messageHandle{}

type ChatService struct {
}

func (is *ChatService) StreamMessages(csi Chat_StreamMessagesServer) error {

	clientUniqueCode := uuid.New().String()
	errch := make(chan error)

	// receive messages - init a go routine
	go receiveFromStream(csi, clientUniqueCode, errch)

	// send messages - init a go routine
	go sendToStream(csi, clientUniqueCode, errch)

	return <-errch

}

// receive messages
func receiveFromStream(csi_ Chat_StreamMessagesServer, clientUniqueCode_ string, errch_ chan error) {

	//implement a loop
	for {
		msg, err := csi_.Recv()
		if err != nil {
			log.Printf("Error in receiving message from client :: %v", err)
			errch_ <- err
		} else {

			messageHandleObject.mu.Lock()

			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        msg.User,
				MessageBody:       msg.Text,
				MessageUniqueCode: uuid.New().String(),
				ClientUniqueCode:  clientUniqueCode_,
			})

			log.Printf("%v", messageHandleObject.MQue[len(messageHandleObject.MQue)-1])

			messageHandleObject.mu.Unlock()

		}
	}
}

// send message
func sendToStream(csi_ Chat_StreamMessagesServer, clientUniqueCode_ string, errch_ chan error) {
	for {

		for {

			time.Sleep(500 * time.Millisecond)

			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			senderUniqueCode := messageHandleObject.MQue[0].ClientUniqueCode
			senderName4Client := messageHandleObject.MQue[0].ClientName
			message4Client := messageHandleObject.MQue[0].MessageBody

			messageHandleObject.mu.Unlock()

			if senderUniqueCode != clientUniqueCode_ {

				err := csi_.Send(&Message{User: senderName4Client, Text: message4Client})

				if err != nil {
					errch_ <- err
				}

				messageHandleObject.mu.Lock()

				if len(messageHandleObject.MQue) > 1 {
					messageHandleObject.MQue = messageHandleObject.MQue[1:]
				} else {
					messageHandleObject.MQue = []messageUnit{}
				}

				messageHandleObject.mu.Unlock()

			}

		}

		time.Sleep(100 * time.Millisecond)
	}
}
