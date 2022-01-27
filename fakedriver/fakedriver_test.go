package fakedriver_test

import (
	"database/sql"
	"testing"

	_ "github.com/nobishino/gocode/fakedriver"
)

func TestDriverRegistration(t *testing.T) {
	t.Log("sql.Open()")
	db, err := sql.Open("fakedriver", "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("db.Ping()")
	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	t.Log("db.Close()")
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}
