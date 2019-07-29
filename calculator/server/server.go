package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	log.Println("FindMaximum function was invoked")
	var maximum int64

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("erro while reading client stream: %v", err)
			return err
		}
		number := req.GetNumber()

		if maximum < number {
			maximum = number
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Result: number,
			})
			if err != nil {
				log.Fatalf("error while sending data to client: %v", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Printf("SquareRoot function was invoked with %v", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	log.Println("CALCULATOR SERVER STARTED")

	lis, err := net.Listen("tcp", "0.0.0.0:18075")
	if err != nil {
		log.Fatalf("Failed to listen on prot 18075: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
