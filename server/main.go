package main

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"log"
	pb "megafon-test/api"
	"megafon-test/server/config"
	"megafon-test/server/db"
	"megafon-test/server/internal"
	"net"
	"os"
)

var (
	port = flag.String("port", "5757", "server port")
)

func main() {
	conf := config.ParseConfig()
	f, err := os.OpenFile(conf.LogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	listener, err := net.Listen("tcp", ":"+(*port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	database, err := db.NewDatabase(conf.DbPath)
	if err != nil {
		log.Fatalf(err.Error())
	}
	pb.RegisterCarPositionServer(grpcServer, &internal.Server{Db: database})
	log.Printf("server listening on %v", *port)
	grpcServer.Serve(listener)
}
