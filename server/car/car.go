package car

import (
	"database/sql"
	"math"
	"megafon-test/server/db"
)

type Car struct {
	Id     int64
	Xcoord int64
	Ycoord int64
}

type CarList struct {
	Cars []Car
}

func (c *CarList) At(i int) db.Storable {
	return &c.Cars[i]
}

func (c *CarList) Len() int {
	return len(c.Cars)
}

func (c *CarList) AppendEmpty() {
	c.Cars = append(c.Cars, Car{0, 0, 0})
}

func (car *Car) StoreSelf(db *sql.DB) error {
	stmt := `INSERT INTO Positions (Id, XCoord, YCoord) VALUES (?, ?, ?) ON CONFLICT(Id) DO UPDATE SET Id = ?, XCoord = ?, YCoord = ?`
	_, err := db.Exec(stmt, car.Id, car.Xcoord, car.Ycoord, car.Id, car.Xcoord, car.Ycoord)
	if err != nil {
		return err
	}
	return nil
}

func (car *Car) GetDist(other Car) float64 {
	return math.Sqrt(float64((car.Xcoord-other.Xcoord)*(car.Xcoord-other.Xcoord) + (car.Ycoord-other.Ycoord)*(car.Ycoord-other.Ycoord)))
}

func (car *Car) SelfScan(row *sql.Rows) error {
	err := row.Scan(&car.Id, &car.Xcoord, &car.Ycoord)
	if err != nil {
		return err
	}
	return nil
}
