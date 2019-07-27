package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vitorhcastro/grpc-master-class/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("---------------------------------------------------------------------------")
	log.Println("Starting Client Streaming RPC")
	fmt.Println()

	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vitor",
				LastName:  "Castro",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "James",
				LastName:  "Bond",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range request {
		log.Printf("Sending request %v", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}

	log.Println(res)

	fmt.Println()
	log.Println("Ending Client Streaming RPC")
	fmt.Println("---------------------------------------------------------------------------")
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("---------------------------------------------------------------------------")
	log.Println("Starting Server Streaming RPC")
	fmt.Println()

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vitor",
			LastName:  "Castro",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading the Stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
	fmt.Println()
	log.Println("Ending Server Streaming RPC")
	fmt.Println("---------------------------------------------------------------------------")
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("---------------------------------------------------------------------------")
	log.Println("Starting Unary RPC")
	fmt.Println()

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vitor",
			LastName:  "Castro",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res)

	fmt.Println()
	log.Println("Ending Unary RPC")
	fmt.Println("---------------------------------------------------------------------------")
}
