// これがモジュールパス(module path)
module github.com/nobishino/gocode

go 1.21.5

// この部分がdependency information
require (
	github.com/google/go-cmp v0.5.7
	github.com/pkg/errors v0.9.1
)

require golang.org/x/time v0.1.0 // indirect
