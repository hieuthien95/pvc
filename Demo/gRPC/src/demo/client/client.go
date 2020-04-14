package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "demo/customer"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

// createCustomer calls the RPC method CreateCustomer of CustomerServer
func createCustomer(client pb.CustomerClient, customer *pb.CustomerRequest) {

	resp, err := client.CreateCustomer(context.Background(), customer)
	if err != nil {
		log.Fatalf("Could not create Customer: %v", err)
	}
	if resp.Success {
		log.Printf("A new Customer has been added with id: %d", resp.Id)
	}
}

// createMultiCustomer calls RPC
func createMultiCustomer(client pb.CustomerClient) {

	stream, err := client.CreateMultiCustomers(context.Background())

	if err != nil {
		log.Fatalf("Could not create Customer: %v", err)
	}
	var arrCus []*pb.CustomerRequest
	customer := &pb.CustomerRequest{
		Id:    3,
		Name:  "CongPV 3",
		Email: "vancong3@gmail.com",
		Phone: "123456789",
		Addresses: []*pb.CustomerRequest_Address{
			{
				Street:            "111F Nguyen Lam",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: false,
			},
			{
				Street:            "111E Nguyen Lam",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: true,
			},
		},
	}
	arrCus = append(arrCus, customer)

	customer = &pb.CustomerRequest{
		Id:    4,
		Name:  "CongPV 4",
		Email: "vancong4@gmail.com",
		Phone: "1234567890",
		Addresses: []*pb.CustomerRequest_Address{
			{
				Street:            "304 To Hien Thanh",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: true,
			},
		},
	}

	arrCus = append(arrCus, customer)

	for _, value := range arrCus {
		if err := stream.Send(value); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, value, err)
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}

	log.Printf("Create success: %v", reply)

}

// getCustomers calls the RPC method GetCustomers of CustomerServer
func getCustomers(client pb.CustomerClient, filter *pb.CustomerFilter) {

	// calling the streaming API
	stream, err := client.GetCustomers(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get customers: %v", err)
	}
	for {
		customer, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetCustomers(_) = _, %v", client, err)
		}
		log.Printf("Customer: %v", customer)
	}

}

func getMultiCustomers(client pb.CustomerClient) {

	stream, err := client.GetMultiCustomers(context.Background())
	if err != nil {
		log.Fatalf("Error on get multi customers: %v", err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			customer, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("%v.GetMultiCustomers(_) = _, %v", client, err)
			}
			log.Printf("Customer: %v", customer)
		}
	}()

	var arrFlilter []*pb.CustomerFilter
	filter := &pb.CustomerFilter{Keyword: "CongPV 1"}
	arrFlilter = append(arrFlilter, filter)

	filter = &pb.CustomerFilter{Keyword: "CongPV 4"}
	arrFlilter = append(arrFlilter, filter)

	fmt.Println(arrFlilter)
	for _, value := range arrFlilter {
		if err = stream.Send(value); err != nil {
			log.Fatalf("Failed to get a customer: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func main() {

	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Creates a new CustomerClient
	client := pb.NewCustomerClient(conn)

	customer := &pb.CustomerRequest{
		Id:    1,
		Name:  "CongPV 1",
		Email: "vancong1@gmail.com",
		Phone: "123456789",
		Addresses: []*pb.CustomerRequest_Address{
			{
				Street:            "111C Nguyen Lam",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: false,
			},
			{
				Street:            "111B Nguyen Lam",
				City:              "TPHCM",
				State:             "TP",
				Zip:               "124",
				IsShippingAddress: true,
			},
		},
	}

	//Create a new customer
	createCustomer(client, customer)

	// Create multi customer
	createMultiCustomer(client)

	// Filter with an empty Keyword
	filter := &pb.CustomerFilter{Keyword: ""}
	getCustomers(client, filter)

	// Filter request multi key
	getMultiCustomers(client)
}
