package main

import (
	"context"
	"log"
	"math/rand"
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		doSum(c, r.Int63()%10000, r.Int63()%10000)
		time.Sleep(time.Duration(r.Int63()) % 10000)
	}
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
