package main

import "testing"

// run with "go test -race", then it fails
func TestRaceDetection(t *testing.T) {
	main()
}

func TestIncorrectSync(t *testing.T) {
	incorrectSync()
}
