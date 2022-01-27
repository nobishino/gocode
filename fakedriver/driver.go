package fakedriver

import (
	"database/sql"
	"database/sql/driver"
	"log"
)

type Driver struct{}

var d Driver

func init() {
	log.Println("sql.Register")
	sql.Register("fakedriver", &d)
}

func (d *Driver) Open(name string) (driver.Conn, error) {
	log.Println("Driver.Open")
	return &Conn{}, nil
}
