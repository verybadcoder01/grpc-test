package main

import (
	"bufio"
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	pb "megafon-test/api"
	"os"
	"strconv"
	"time"
)

var (
	port = flag.String("port", "5757", "server port")
)

func ReadInt(reader *bufio.Reader) int64 {
	data, _ := reader.ReadString('\n')
	data = data[:len(data)-1]
	res, _ := strconv.Atoi(data)
	return int64(res)
}

func ReadFloat(reader *bufio.Reader) float32 {
	data, _ := reader.ReadString('\n')
	data = data[:len(data)-1]
	res, _ := strconv.ParseFloat(data, 32)
	return float32(res)
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("127.0.0.1"+":"+(*port), opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewCarPositionClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	var action string
	reader := bufio.NewReader(os.Stdin)
	for true {
		action, err = reader.ReadString('\n')
		action = action[:len(action)-1]
		if action == "stop" {
			break
		}
		if action == "store" {
			car := &pb.Car{Id: &pb.Id{CarId: ReadInt(reader)}, Point: &pb.Coords{Xcoord: ReadInt(reader), Ycoord: ReadInt(reader)}}
			resp, err := client.Store(ctx, car)
			log.Print(resp)
			if err != nil {
				log.Println(err.Error())
			}
		} else if action == "retrieve" {
			resp, err := client.Retrieve(ctx, &pb.Id{CarId: ReadInt(reader)})
			log.Print(resp)
			if err != nil {
				log.Println(err.Error())
			}
		} else if action == "circle" {
			resp, err := client.Neighbors(ctx, &pb.Circle{Point: &pb.Coords{Xcoord: ReadInt(reader), Ycoord: ReadInt(reader)}, Area: ReadFloat(reader)})
			log.Print(resp)
			if err != nil {
				log.Println(err.Error())
			}
		}
		if err != nil {
			log.Print(err.Error())
		}
	}
}
