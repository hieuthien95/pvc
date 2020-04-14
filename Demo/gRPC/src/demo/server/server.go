package main

import (
	"context"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	pb "demo/customer"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement customer.CustomerServer.
type server struct {
	savedCustomers []*pb.CustomerRequest
	mu             sync.Mutex // protects routeNotes
}

// Create multi customer
func (s *server) CreateMultiCustomers(stream pb.Customer_CreateMultiCustomersServer) error {

	var countCreate int32
	for {
		value, err := stream.Recv()

		if err == io.EOF { // Hoan thanh nhan request
			return stream.SendAndClose(&pb.CustomerResponse{
				Id:      countCreate,
				Success: true,
			})
		}
		s.savedCustomers = append(s.savedCustomers, value)

		countCreate++
	}
}

// CreateCustomer creates a new Customer
func (s *server) CreateCustomer(ctx context.Context, in *pb.CustomerRequest) (*pb.CustomerResponse, error) {

	s.savedCustomers = append(s.savedCustomers, in)
	return &pb.CustomerResponse{
		Id:      in.Id,
		Success: true,
	}, nil
}

func (s *server) GetMultiCustomers(stream pb.Customer_GetMultiCustomersServer) error {

	for {
		in, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		filter := pb.CustomerFilter{Keyword: in.Keyword}
		for _, customer := range s.savedCustomers {

			if filter.Keyword != "" {
				if !strings.Contains(customer.Name, filter.Keyword) {
					continue
				}
			}
			if err := stream.Send(customer); err != nil {
				return err
			}
		}
	}
}

// GetCustomers returns all customers by given filter
func (s *server) GetCustomers(filter *pb.CustomerFilter, stream pb.Customer_GetCustomersServer) error {

	for _, customer := range s.savedCustomers {
		if filter.Keyword != "" {
			if !strings.Contains(customer.Name, filter.Keyword) {
				continue
			}
		}
		if err := stream.Send(customer); err != nil {
			return err
		}
	}
	return nil
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterCustomerServer(s, &server{})
	s.Serve(lis)
}
