package main

import (
	"context"
	"fmt"
	"gitlab.eng.vmware.com/nidhig1/go_grpc_assignment/calculator/calcApi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

func main() {
	fmt.Println("Hello I am Client")
	cc, error := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if error != nil {
		fmt.Println("Error while connecting to server")
	}
	defer cc.Close()
	c := calcApi.NewCalcServiceClient(cc)
	//fmt.Println(" Unary RPC for calculating sum")
	//calculateSum(c)
	//fmt.Println(" Server Streaming RPC for multiplication table")
	//calculateMultiplicationTable(c)
	//fmt.Println(" Client Streaming RPC for maximum number")
	//calculateMaximumNumber(c)
	fmt.Println(" Bidirectional Streaming RPC for average")
	calculateAverage(c)
}

func calculateSum(c calcApi.CalcServiceClient) {
	request := &calcApi.SumRequest{
		Input: &calcApi.InputNumbers{
			FirstNum:  56,
			SecondNum: 72,
		},
	}
	fmt.Printf("Sending %v to server as inputs \n", request)
	response, err := c.Sum(context.Background(), request)
	if err != nil {
		fmt.Println("Error while calling API")
	}
	fmt.Println("Response From API: Sum of Two number is=", response.Response)
}

func calculateMultiplicationTable(c calcApi.CalcServiceClient) {
	request := &calcApi.MultiplicationTableRequest{
		Number: 9,
	}
	fmt.Printf("Sending %v to server as input \n", request)
	resStream, err := c.MultiplicationTable(context.Background(), request)
	if err != nil {
		fmt.Println("Error while sending request")
	}
	for {
		message, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while receiving data from server")
		}
		fmt.Println("Response from server streaming", message.GetResponse())
	}
}

func calculateMaximumNumber(c calcApi.CalcServiceClient) {
	request := []*calcApi.MaximumNumberRequest{
		&calcApi.MaximumNumberRequest{
			Number: 68,
		},
		&calcApi.MaximumNumberRequest{
			Number: 45,
		},
		&calcApi.MaximumNumberRequest{
			Number: 90,
		},
		&calcApi.MaximumNumberRequest{
			Number: 99,
		},
		&calcApi.MaximumNumberRequest{
			Number: 87,
		},
	}

	stream, err := c.Maximum(context.Background())
	if err != nil {
		fmt.Println("Error while calling Maximum Number RPC")
	}
	for _, req := range request {
		fmt.Printf("Sending %v to server as input \n", req)
		stream.Send(req)
	}
	response, er := stream.CloseAndRecv()
	if er != nil {
		fmt.Println("Error while receiving response from Maximum rpc")
	}
	fmt.Println("Response from Maximum-The Maximum No is:", response.GetResponse())
}

func calculateAverage(c calcApi.CalcServiceClient) {
	request := []*calcApi.AverageRequest{
		&calcApi.AverageRequest{
			Number: 68,
		},
		&calcApi.AverageRequest{
			Number: 45,
		},
		&calcApi.AverageRequest{
			Number: 90,
		},
		&calcApi.AverageRequest{
			Number: 99,
		},
		&calcApi.AverageRequest{
			Number: 87,
		},
	}

	stream, err := c.Average(context.Background())
	if err != nil {
		fmt.Println("Error while creating stream")
	}
	endChannel := make(chan struct{})
	go func() {
		for _, req := range request {
			fmt.Printf("Sending %v to Server\n", req)
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
				fmt.Println("Error while receiving response from average")
				break
			}
			fmt.Println("Average Received", res.GetResponse())
		}
		close(endChannel)
	}()
	<-endChannel
}
