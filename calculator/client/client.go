package main

import (
	"context"
	"log"

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
	doSum(c, 1, 2)
}

func doSum(c calculatorpb.CalculatorServiceClient, first int64, second int64) {
	log.Printf("Starting to do a Sum RPC of values %v and %v\n", first, second)
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
}
