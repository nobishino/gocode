package fakedriver

import (
	"context"
	"database/sql/driver"
	"log"
)

type Connector struct {
}

var _ driver.Connector = &Connector{}

func (c *Connector) Connect(context.Context) (driver.Conn, error) {
	log.Println("Connector.Connect")
	return &Conn{}, nil
}

func (c *Connector) Driver() driver.Driver {
	log.Println("Connector.Driver")
	return &d
}
