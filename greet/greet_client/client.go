package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vitorhcastro/grpc-master-class/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	// doUnary(c)

	// doServerStreaming(c)

	// doClientStreaming(c)

	// doBiDiStreaming(c)

	doUnaryWithDeadlines(c, 1*time.Second) // should timeout
	doUnaryWithDeadlines(c, 5*time.Second) // should complete
}

func doUnaryWithDeadlines(c greetpb.GreetServiceClient, timeout time.Duration) {
	fmt.Println("---------------------------------------------------------------------------")
	log.Println("Starting UnaryWithDeadline RPC")
	fmt.Println()

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Vitor",
			LastName:  "Castro",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Deadline was exceeded.")
			} else {
				fmt.Printf("unexpected error: %v", err)
			}
		} else {
			log.Fatalf("error while calling Greet RPC: %v", err)
		}

		fmt.Println()
		log.Println("Ending UnaryWithDeadline RPC")
		fmt.Println("---------------------------------------------------------------------------")
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res)

	fmt.Println()
	log.Println("Ending UnaryWithDeadline RPC")
	fmt.Println("---------------------------------------------------------------------------")
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("---------------------------------------------------------------------------")
	log.Println("Starting BiDI Streaming RPC")
	fmt.Println()

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Vitor",
				LastName:  "Castro",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "James",
				LastName:  "Bond",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v", err)
		return
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range requests {
			log.Printf("Sending message %v", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving response from GreetEveryone: %v", err)
			}
			log.Printf("Received message %v", res)
		}
		close(waitc)
	}()

	<-waitc

	fmt.Println()
	log.Println("Ending BiDI Streaming RPC")
	fmt.Println("---------------------------------------------------------------------------")
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
