package server

import (
	"context"
	"encoding/json"
	pb "github.com/AnyKey/userslike/grpcsrv/like"
	"github.com/AnyKey/userslike/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

type server struct {
	Repo repository.Repository
	pb.UnimplementedSubSrvServer
}

func (s *server) SetLike(ctx context.Context, in *pb.LikeRequest) (*pb.LikeReply, error) {
	myToken := in.GetJwt()
	claims := jwt.MapClaims{}
	res, err := jwt.ParseWithClaims(myToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	if !res.Valid {
		return nil, errors.New("unauthorized")
	}
	user := (claims["name"]).(string)
	err = s.Repo.SetLike(in.GetName(), in.GetArtist(), user)
	if err != nil {
		return nil, err
	}
	return &pb.LikeReply{Message: "OK "}, nil
}
func (s *server) GetLike(ctx context.Context, in *pb.TrackRequest) (*pb.TrackReply, error) {
	myToken := in.GetJwt()
	claims := jwt.MapClaims{}
	res, err := jwt.ParseWithClaims(myToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("key"), nil
	})
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	if !res.Valid {
		return nil, errors.New("unauthorized")
	}
	result, lk, err := s.Repo.GetTracks(in.GetName(), in.GetArtist())
	if err != nil {
		return nil, err
	}
	bytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return &pb.TrackReply{
		Name:      in.GetName(),
		Artist:    in.GetArtist(),
		User:      bytes,
		LikeCount: *lk,
	}, nil
}

func Run(repo repository.Repository) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSubSrvServer(s, &server{Repo: repo})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
