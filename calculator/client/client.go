package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/vitorhcastro/grpc-master-class/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Hello, I'm a client :D")

	cc, err := grpc.Dial("localhost:18075", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doSum(c, 1, 2)

	// doPrimeDecomposition(c, 120)

	doCalculateAverage(c, 1, 2, 3, 4)
}

func doCalculateAverage(c calculatorpb.CalculatorServiceClient, numbers ...int64) {
	fmt.Println("---------------------------------------------------------------------------")
	log.Println("Starting CalculateAverage RPC")
	fmt.Println()

	stream, err := c.CalculateAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling CalculateAverage: %v", err)
	}

	for _, number := range numbers {
		req := &calculatorpb.CalculateAverageRequest{
			Number: number,
		}
		log.Printf("Sending request %v", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from CalculateAverage: %v", err)
	}

	log.Println(res)

	fmt.Println()
	log.Println("Ending CalculateAverage RPC")
	fmt.Println("---------------------------------------------------------------------------")
}

func doPrimeDecomposition(c calculatorpb.CalculatorServiceClient, number int64) {
	fmt.Println("---------------------------------------------------------------------------")
	log.Printf("Starting to do a PrimeNumberDecompositon RPC of number %v\n", number)
	fmt.Println()

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: number,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	i := 0
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading the Stream: %v", err)
		}
		i++
		log.Printf("Response number %v from PrimeNumberDecomposition: %v", i, msg.GetResult())
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println()
	log.Println("Ending Server Streaming RPC")
	fmt.Println("---------------------------------------------------------------------------")
}

func doSum(c calculatorpb.CalculatorServiceClient, first int64, second int64) {
	fmt.Println("---------------------------------------------------------------------------")
	log.Printf("Starting to do a Sum RPC of values %v and %v\n", first, second)
	fmt.Println()

	req := &calculatorpb.CalculatorRequest{
		Operation: &calculatorpb.Operation{
			FirstNumber:  first,
			SecondNumber: second,
		},
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("%v + %v = %v", first, second, res.Result)
	log.Printf("Response from Sum: %v", res.Result)

	fmt.Println()
	log.Println("Ending Sum RPC")
	fmt.Println("---------------------------------------------------------------------------")
}
