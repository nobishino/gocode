package main

import "testing"

func TestRule30(t *testing.T) {
	cases := []struct {
		left  bool
		self  bool
		right bool
		want  bool
	}{}
	for _, tc := range cases {
		got := rule30(tc.left, tc.self, tc.right)
		if got != tc.want {
			t.Errorf("want %t, but got %t for input (%t,%t,%t)", tc.want, got, tc.left, tc.self, tc.right)
		}
	}
}
