package fakedriver

import (
	"database/sql"
	"database/sql/driver"
	"log"
)

type Driver struct{}

type DriverContext struct {
	*Driver
}

var d Driver
var dc DriverContext

var _ driver.Driver = &d

var _ driver.DriverContext = &dc

func init() {
	log.Println("sql.Register")
	sql.Register("fakedriver", &d)
	sql.Register("fakedrivercontext", &dc)
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	log.Printf("Driver.Open(%s)\n", name)
	return &Conn{}, nil
}

func (d *DriverContext) OpenConnector(name string) (driver.Connector, error) {
	log.Printf("Driver.OpenConnector(%s)\n", name)
	return &Connector{}, nil
}
