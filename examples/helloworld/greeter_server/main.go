/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
func (s *server) SayBye(ctx context.Context, in *pb.ByeRequest) (*pb.ByeReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.ByeReply{Message: "Hello " + in.GetName() + " " + in.GetInfo()}, nil
}
func (s *server) GetAge(ctx context.Context, in *pb.HelloRequest) (*pb.AgeReply, error) {
	log.Printf("Received: %v", in.GetAge())
	return &pb.AgeReply{Age: in.GetAge()}, nil
}

// stream Demo
func (s *server) KeepSayHello(in *pb.HelloRequest, stream pb.Greeter_KeepSayHelloServer) error {
	age := int(in.GetAge())
	for i := 0; i < age; i++ {
		stream.Send(&pb.HelloReply{Message: "hello" + strconv.Itoa(i)})
	}
	return nil
}

func (s *server) GetAllCongras(stream pb.Greeter_GetAllCongrasServer) error {
	congras := []string{}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&pb.HelloReply{Message: "Hello " + strings.Join(congras, ",")})
			log.Printf("DONE")
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		congras = append(congras, in.GetName())
	}

}

func (s *server) KeepReply(stream pb.Greeter_KeepReplyServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		info := in.GetName() + strconv.Itoa(int(in.GetAge()))
		stream.Send(&pb.HelloReply{Message: info})
	}
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
