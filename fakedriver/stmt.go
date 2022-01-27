package fakedriver

import (
	"database/sql/driver"
	"log"
)

type Stmt struct{}

var _ driver.Stmt = &Stmt{}

func (stmt *Stmt) Close() error {
	log.Println("Stmt.Close")
	return nil
}

func (stmt *Stmt) NumInput() int {
	log.Println("Stmt.NumInput")
	return 0
}

func (stmt *Stmt) Exec(args []driver.Value) (driver.Result, error) {
	log.Println("Stmt.Exec")
	return &Result{}, nil
}

func (stmt *Stmt) Query(args []driver.Value) (driver.Rows, error) {
	log.Println("Stmt.Query")
	return &Rows{}, nil

}
