package main

import "testing"

func TestFileArgumentIsTheFirstNonFlag(t *testing.T) {
	cases := map[string][]string{
		"":        {},
		"a.cuppa": {"a.cuppa"},
		"b.cuppa": {"--flag", "b.cuppa", "c.cuppa"},
	}
	for want, args := range cases {
		if got := fileArgument(args); got != want {
			t.Errorf("%v: got %q, want %q", args, got, want)
		}
	}
}
