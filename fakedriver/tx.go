package fakedriver

import (
	"database/sql/driver"
	"log"
)

type Tx struct{}

var _ driver.Tx = &Tx{}

func (tx *Tx) Commit() error {
	log.Println("Tx.Commit")
	return nil
}

func (tx *Tx) Rollback() error {
	log.Println("Tx.Rollback")
	return nil
}
