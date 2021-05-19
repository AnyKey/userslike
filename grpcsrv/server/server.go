package server

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "users_like/grpcsrv/like"
	"users_like/repository"
)

const (
	port = ":50051"
)

type server struct {
	Repo repository.Repository
	pb.UnimplementedGreeterServer
}

func (s *server) SetLike(ctx context.Context, in *pb.LikeRequest) (*pb.LikeReply, error) {
	log.Printf("Track: %v, Artist: %v, JWT: %v", in.GetName(), in.GetArtist(), in.GetJwt())

	return &pb.LikeReply{Message: "OK "}, nil
}
func (s *server) GetLike(ctx context.Context, in *pb.TrackRequest) (*pb.TrackReply, error) {
	log.Printf("Track: %v, Artist: %v", in.GetName(), in.GetArtist())
	s.Repo.GetTracks(in.GetName(),in.GetArtist())
	return &pb.TrackReply{
		Track:    "Sound",
		Username: "Dude",
	}, nil
}

func Run( repo repository.Repository) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{Repo: repo})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
