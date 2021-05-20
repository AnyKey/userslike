package main

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"users_like/repository"

	"google.golang.org/grpc"
	pb "users_like/grpcsrv/like"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSubSrvClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SetLike(ctx, &pb.LikeRequest{
		Name:   "New Track",
		Artist: "NewArtist",
		Jwt:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjE1MDgyNzksIm5hbWUiOiJOZXdVc2VyIiwicm9vdCI6dHJ1ZX0.hSQyR21DF3jMpjzIl4sPCO5w3gZj-peG6xdxqlJcUc4",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.GetMessage())
	var re *pb.TrackReply
	re, err = c.GetLike(ctx, &pb.TrackRequest{Name: "New Track",
		Artist: "NewArtist",
		Jwt:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjE1MDgyNzksIm5hbWUiOiJOZXdVc2VyIiwicm9vdCI6dHJ1ZX0.hSQyR21DF3jMpjzIl4sPCO5w3gZj-peG6xdxqlJcUc4",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	var likes []repository.LikeSelect
	users := json.Unmarshal(re.GetUser(), &likes)
	log.Println(re.GetName(), re.GetArtist(), users)
}
