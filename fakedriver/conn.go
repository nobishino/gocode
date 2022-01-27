package fakedriver

import (
	"database/sql/driver"
	"log"
)

type Conn struct{}

var _ driver.Conn = &Conn{}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	log.Println("Conn.Prepare")
	return &Stmt{}, nil
}

func (c *Conn) Close() error {
	log.Println("Conn.Close")
	return nil
}

func (c *Conn) Begin() (driver.Tx, error) {
	log.Println("Conn.Begin")
	return &Tx{}, nil
}
