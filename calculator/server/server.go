package main

import (
	"context"
	"log"
	"net"

	"github.com/vitorhcastro/grpc-master-class/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	log.Printf("Sum function was invoked with %v", req)
	firstNumber := req.GetOperation().GetFirstNumber()
	secondNumber := req.GetOperation().GetSecondNumber()
	result := firstNumber + secondNumber
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	log.Println("CALCULATOR SERVER STARTED")

	lis, err := net.Listen("tcp", "0.0.0.0:18075")
	if err != nil {
		log.Fatalf("Failed to listen on prot 18075: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
