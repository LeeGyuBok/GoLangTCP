package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"

	pb "golangTCP/messages"

	"google.golang.org/protobuf/proto"

	mg "golangTCP/packages/manager"
)

var playerManager *mg.PlayerManager

func main() {

	playerManager = mg.NewPlayerManager()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()
	fmt.Println("Server is listening on :8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		// 메시지 길이를 먼저 읽습니다 (4바이트)
		lengthBuf := make([]byte, 4)
		_, err := conn.Read(lengthBuf)
		if err != nil {
			log.Printf("Failed to read message length: %v", err)
			return
		}
		length := binary.LittleEndian.Uint32(lengthBuf)

		// 메시지 본문을 읽습니다
		messageBuf := make([]byte, length)
		_, err = conn.Read(messageBuf)
		if err != nil {
			log.Printf("Failed to read message body: %v", err)
			return
		}

		// Protocol Buffers 메시지를 파싱합니다
		message := &pb.GameMessage{}
		err = proto.Unmarshal(messageBuf, message)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		// 메시지 처리
		processMessage(message)

		// 응답 메시지 생성 및 전송 (예: 에코)
		response, err := proto.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			continue
		}

		// 메시지 길이를 먼저 보냅니다
		binary.BigEndian.PutUint32(lengthBuf, uint32(len(response)))
		conn.Write(lengthBuf)

		// 메시지 본문을 보냅니다
		conn.Write(response)
	}
}

func processMessage(message *pb.GameMessage) {
	switch msg := message.Message.(type) {
	case *pb.GameMessage_PlayerPosition:
		pos := msg.PlayerPosition
		fmt.Println("Position : ", pos.X, pos.Y, pos.Z)
	case *pb.GameMessage_Chat:
	case *pb.GameMessage_Login:
		playerId := msg.Login.PlayerId
		fmt.Println(playerId)
		playerManager.AddPlayer(playerId, 0)
	default:
		panic(fmt.Sprintf("unexpected messages.isGameMessage_Message: %#v", msg))
	}
}
