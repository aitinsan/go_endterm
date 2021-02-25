package main

import (
	"endterm/calculator/calculatorpb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{}

func (s server) Max(stream calculatorpb.CalculatorService_MaxServer) error {
	fmt.Println("Max service invoked")
	max := int32(0)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {return nil}
		if err != nil {log.Fatalf("Error when receiving stream: %v", err)}
		n := req.GetNumber()
		if n > max {
			max = n
			err := stream.Send(&calculatorpb.MaxResponse{Result: max})
			if err != nil {log.Fatalf("Error when sending stream: %v", err)}
		}
	}
}

func (s server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	fmt.Println("Average service invoked")
	nums := int32(0)
	count := int32(0)

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			res := float32(nums)/float32(count)
			err := stream.SendAndClose(&calculatorpb.AverageResponse{Result: res})
			if err != nil {log.Fatalf("Error when sending response: %v", err)}
			return nil
		}
		if err != nil {
			log.Fatalf("Error while streaming: %v", err)
		}
		nums += req.GetNumber()
		count++
	}
}

func (s server) PrimeNumberDecomposition(
	req *calculatorpb.PrimeNumberRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer,
) error {
	fmt.Println("Prime Number Decomposition Service Invoked")
	n := req.GetNumber()
	primeNumberDecomposition(n, stream)
	return nil
}

func (s server) Sum(_ context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Println("Sum Service Invoked")
	fn := req.GetFirstNum()
	sn := req.GetSecondNum()
	result := fn + sn
	return &calculatorpb.SumResponse{Result: result}, nil
}

func main() {
	fmt.Println("Starting up Server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}

func primeNumberDecomposition(n int32, steam calculatorpb.CalculatorService_PrimeNumberDecompositionServer) {
	k := int32(2)
	for n > 1 {
		if n%k == 0 {
			err := steam.Send(&calculatorpb.PrimeNumberResponse{Result: k})
			if err != nil {
				log.Fatalf("Error while attempting to send stream: %v", err)
			}
			n /= k
		} else {
			k++
		}
	}
}
