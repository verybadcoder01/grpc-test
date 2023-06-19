package internal

import (
	"context"
	"log"
	pb "megafon-test/api"
	"megafon-test/server/car"
	"megafon-test/server/db"
)

type Server struct {
	Db *db.Database
}

func (s *Server) Store(c context.Context, request *pb.Car) (*pb.Status, error) {
	req := &car.Car{Id: request.Id.GetCarId(), Xcoord: request.Point.GetXcoord(), Ycoord: request.Point.GetYcoord()}
	err := s.Db.Store(req)
	var status = "success"
	if err != nil {
		status = "fail"
		log.Println(err)
		return &pb.Status{Status: status}, err
	}
	return &pb.Status{Status: status}, nil
}

func (s *Server) Retrieve(c context.Context, request *pb.Id) (*pb.Coords, error) {
	var res car.Car
	err := s.Db.GetStorable(request.CarId, &res)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &pb.Coords{Xcoord: res.Xcoord, Ycoord: res.Ycoord}, nil
}

func (s *Server) Neighbors(c context.Context, request *pb.Circle) (*pb.Cars, error) {
	var cars car.CarList
	err := s.Db.GetAll(&cars)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var res []*pb.Car
	base := car.Car{Xcoord: request.Point.GetXcoord(), Ycoord: request.Point.GetYcoord()}
	for _, value := range cars.Cars {
		cur := car.Car{Id: value.Id, Xcoord: value.Xcoord, Ycoord: value.Ycoord}
		if float32(base.GetDist(cur)) <= request.Area {
			res = append(res, &pb.Car{Id: &pb.Id{CarId: cur.Id}, Point: &pb.Coords{Xcoord: cur.Xcoord, Ycoord: cur.Ycoord}})
		}
	}
	return &pb.Cars{Cars: res}, nil
}
