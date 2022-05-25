package main

import (
	"context"
	"log"

	pb "github.com/supertokens/supertokens-golang/examples/with-grpc/hatmaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr string = "localhost:8080"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to connect to server: " + err.Error())
	}

	defer conn.Close()

	c := pb.NewHatmakerClient(conn)

	createHat(c)
}

func createHat(c pb.HatmakerClient) (*pb.Hat, error) {
	log.Println("---createHat was involked---")
	size := &pb.Size{
		Inches: 13,
	}
	res, err := c.MakeHat(context.Background(), size)
	if err != nil {
		log.Println("Unexpected error: ", err.Error())
		return nil, err
	}

	log.Println("Hat has been created", res)
	return res, nil
}
