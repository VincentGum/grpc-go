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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	defaultName = "world"
	defaultInfo = "see u later"
	defaultAge  = 11
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
	info = flag.String("info", defaultInfo, "Info to bye")
	age  = flag.Int("age", defaultAge, "age")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())

	r2, err := c.SayBye(ctx, &pb.ByeRequest{Name: "john", Info: "hope to see u again"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r2.GetMessage())

	r3, err := c.GetAge(ctx, &pb.HelloRequest{Name: "john", Age: defaultAge})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %v", r3.GetAge())

	stream, err := c.KeepSayHello(context.Background(), &pb.HelloRequest{Name: "john", Age: defaultAge})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			log.Println("DONE")
			break
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
		}
		log.Printf("Greeting: %s", reply.Message)
	}

	allCongrasStream, err := c.GetAllCongras(context.Background())
	names := []string{"john", "tom", "jack"}
	for i := 0; i < 3; i++ {
		if err != nil {
			log.Printf("failed to call: %v", err)
			break
		}
		allCongrasStream.Send(&pb.HelloRequest{Name: names[i]})
	}
	allCongrasReply, err := allCongrasStream.CloseAndRecv()
	if err != nil {
		log.Printf("failed to recv: %v", err)
	}
	log.Printf("Greeting: %s", allCongrasReply.Message)

	KeepReplyStream, err := c.KeepReply(context.Background())
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	j := 0
	for {
		err := KeepReplyStream.Send(&pb.HelloRequest{Name: names[j], Age: int64(j)})
		if err != nil {
			log.Printf("failed to send: %v", err)
			break
		}
		KeepReplyStreamReply, err := KeepReplyStream.Recv()
		if err != nil {
			log.Printf("failed to recv: %v", err)
			break
		}
		log.Printf("Greeting: %s", KeepReplyStreamReply.Message)
		j++
		if j == 3 {
			break
		}
	}
}
