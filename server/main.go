package main

import (
	"flag"
	"google.golang.org/grpc"
	"log"
	pb "megafon-test/api"
	"megafon-test/logger"
	"megafon-test/server/config"
	"megafon-test/server/db"
	"megafon-test/server/internal"
	"net"
)

var (
	port = flag.String("port", "5757", "server port")
)

func main() {
	conf := config.ParseConfig()
	logger.SetupLogging(conf.LogPath, 32, 28, 1)
	listener, err := net.Listen("tcp", ":"+(*port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	database, err := db.NewDatabase(&db.DatabaseConfig{Driver: conf.DbDriver, FilePath: conf.DbPath, DSN: conf.DSN})
	if err != nil {
		log.Fatalf(err.Error())
	}
	pb.RegisterCarPositionServer(grpcServer, &internal.Server{Db: database})
	log.Printf("server listening on %v", *port)
	grpcServer.Serve(listener)
}
