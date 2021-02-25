package main

import (
	"endterm/calculator/calculatorpb"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// callSum(c)
	callPrimeNumber(c)
	callAverage(c)
	//callMax(c)
}



func callPrimeNumber(c calculatorpb.CalculatorServiceClient) {
	n := int32(120)
	req := calculatorpb.PrimeNumberRequest{Number: n}
	rs := make([]int32, 0)

	fmt.Println("Number:", n)

	stream, err := c.PrimeNumberDecomposition(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while calling Prime Num: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("Error while streaming: %v", err)
		}
		rs = append(rs, resp.GetResult())
		fmt.Println("Number decomposes to:", resp.GetResult())
	}
	if len(rs) == 1 {
		fmt.Println(n, "is prime")
	} else {
		fmt.Println(n, "is not prime")
	}
	defer fmt.Println(n, "is equal to the following numbers multiplied together", rs)
}

func callAverage(c calculatorpb.CalculatorServiceClient) {
	nums := []int32{1, 2, 3, 4}
	stream, err := c.Average(context.Background())
	if err != nil {log.Fatalf("Error opening stream to Average: %v", err)}

	for _, num := range nums {
		fmt.Println("Sending", num)
		err := stream.Send(&calculatorpb.AverageRequest{Number: num})
		if err != nil {log.Fatalf("Error streaming data: %v", err)}
		time.Sleep(time.Second)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {log.Fatalf("Error receiving response: %v", err)}

	fmt.Println("Average of numbers is", resp.GetResult())
}


