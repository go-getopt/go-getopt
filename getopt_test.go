package getopt_test

import (
	"testing"

	. "gitee.com/go-getopt/go-getopt"
)

func TestGetOpt(t *testing.T) {
	t.Log(GetOpt([]string{"ls", "-al", "/etc"}, "al", "all"))
	t.Log("Arguments:", Args)
	t.Log("Program name:", Args[0])
	for loop := true; loop; {
		arg1 := Get(1)
		switch arg1 {
		case "-a":
			Shift(1)
		case "-l":
			Shift(1)
		case "--":
			Shift(1)
			loop = false
		default:
			t.Error("Wrong argument'", arg1, "'")
			t.FailNow()
		}
	}
	t.Log("Positional args:", Args[1:])
}
