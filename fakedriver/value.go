package fakedriver

import "database/sql/driver"

type Value struct{}

var _ driver.Value = &Value{}
