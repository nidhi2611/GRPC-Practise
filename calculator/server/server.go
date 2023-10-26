package main

import (
	"context"
	"fmt"
	"gitlab.eng.vmware.com/nidhig1/go_grpc_assignment/calculator/calcApi"
	"google.golang.org/grpc"
	"io"
	"net"
	"strconv"
)

type server struct {
}

func (*server) Average(stream calcApi.CalcService_AverageServer) error {
	var average int64
	var count int
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Println("Error while reading data from Average request")
			return err
		}
		number := request.GetNumber()
		count += 1
		average = average + number
		er := stream.Send(&calcApi.AverageResponse{
			Response: strconv.Itoa(int(average) / count),
		})
		if er != nil {
			fmt.Println("Error while sending data to client")
			return er
		}
	}
}

func (*server) Maximum(stream calcApi.CalcService_MaximumServer) error {
	var max int64
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calcApi.MaximumNumberResponse{
				Response: strconv.Itoa(int(max)),
			})
		}
		if err != nil {
			fmt.Println(" Error while reading request from Maximum API")
		}
		number := response.GetNumber()
		if number > max {
			max = number
		}
	}
}

func (*server) Sum(ctx context.Context, request *calcApi.SumRequest) (*calcApi.SumResponse, error) {
	firstInput := request.GetInput().GetFirstNum()
	secondInput := request.GetInput().GetSecondNum()
	response := &calcApi.SumResponse{
		Response: firstInput + secondInput,
	}
	return response, nil
}

func (*server) MultiplicationTable(request *calcApi.MultiplicationTableRequest, stream calcApi.CalcService_MultiplicationTableServer) error {
	number := request.GetNumber()
	var i int64
	for i = 1; i <= 10; i++ {
		result := strconv.Itoa(int(number)) + " X " + strconv.Itoa(int(i)) + " = " + strconv.Itoa(int(number*i))
		response := &calcApi.MultiplicationTableResponse{
			Response: result,
		}
		stream.Send(response)
	}
	return nil
}

func main() {
	fmt.Println("Starting GRPC server")
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		fmt.Println("Error while starting server")
	}
	s := grpc.NewServer()
	calcApi.RegisterCalcServiceServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		fmt.Println("Failed to serve")
	}
}
