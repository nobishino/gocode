package fakedriver

import (
	"database/sql/driver"
	"io"
	"log"
)

type Rows struct{}

var _ driver.Rows = &Rows{}

func (r *Rows) Columns() []string {
	log.Println("Rows.Columns")
	return []string{}
}

func (r *Rows) Close() error {
	log.Println("Rows.Close")
	return nil
}
func (r *Rows) Next(dest []driver.Value) error {
	log.Println("Rows.Next")
	return io.EOF
}
