package fakedriver

import (
	"database/sql/driver"
	"log"
)

type Result struct{}

var _ driver.Result = &Result{}

func (r *Result) LastInsertId() (int64, error) {
	log.Println("Result.LastInsertId")
	return 0, nil
}

func (r *Result) RowsAffected() (int64, error) {
	log.Println("Result.RowsAffected")
	return 0, nil
}
