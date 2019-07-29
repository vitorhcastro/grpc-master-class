package main

import (
	"context"
	"io"
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

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("PrimeNumberDecomposition function was invoked with %v", req)
	number := req.GetNumber()
	var k int64 = 2
	for number > 1 {
		if number%k == 0 {
			number /= k
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: k,
			}
			stream.Send(res)
		} else {
			k++
		}
	}
	return nil
}
func (*server) CalculateAverage(stream calculatorpb.CalculatorService_CalculateAverageServer) error {
	log.Println("CalculateAverage function was invoked")
	result := 0.0
	count := 0.0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.CalculateAverageResponse{
				Result: result / count,
			})
		}
		if err != nil {
			log.Fatalf("error while receiving stream: %v", err)
		}

		result += float64(req.GetNumber())
		count++
	}
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
