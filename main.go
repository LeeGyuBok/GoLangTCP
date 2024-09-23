package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "golangtcp/messages"

	"google.golang.org/grpc"
)

type gameServer struct {
	pb.UnimplementedGameServiceServer
	mu      sync.Mutex
	players map[int32]*pb.PlayerPosition
}

func (s *gameServer) UpdatePosition(ctx context.Context, pos *pb.PlayerPosition) (*pb.GameState, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.players[pos.PlayerId] = pos

	gameState := &pb.GameState{
		Players: make([]*pb.PlayerPosition, 0, len(s.players)),
	}
	for _, player := range s.players {
		gameState.Players = append(gameState.Players, player)
	}

	log.Printf("Updated position for player %d: (%f, %f, %f)", pos.PlayerId, pos.X, pos.Y, pos.Z)
	return gameState, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGameServiceServer(s, &gameServer{
		players: make(map[int32]*pb.PlayerPosition),
	})

	fmt.Println("Game server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
