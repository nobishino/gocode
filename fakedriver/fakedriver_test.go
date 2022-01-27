package fakedriver_test

import (
	"database/sql"
	"testing"

	_ "github.com/nobishino/gocode/fakedriver"
)

/*
=== RUN   TestDriverRegistration
    fakedriver_test.go:11: sql.Open()
    fakedriver_test.go:17: db.Ping()
2022/01/27 23:04:30 Driver.Open()
    fakedriver_test.go:22: db.Close()
2022/01/27 23:04:30 Conn.Close
--- PASS: TestDriverRegistration (0.00s)
*/
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

/*
=== RUN   TestDriverContextRegistration
    fakedriver_test.go:29: sql.Open()
2022/01/27 23:04:30 Driver.OpenConnector()
    fakedriver_test.go:35: db.Ping()
2022/01/27 23:04:30 Connector.Connect
    fakedriver_test.go:40: db.Close()
2022/01/27 23:04:30 Conn.Close
*/
func TestDriverContextRegistration(t *testing.T) {
	t.Log("sql.Open()")
	db, err := sql.Open("fakedrivercontext", "")
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
