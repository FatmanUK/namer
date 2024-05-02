package main

import (
	"os"
	"fmt"
)

const (
	LL_ERROR = uint16(0)
	LL_WARN  = uint16(1)
	LL_INFO  = uint16(2)
	LL_DEBUG = uint16(3)
)

var prefix_codes = [...]string{"[x.x]", "[O.o]", "[o.o]", "[-.-]"}

type View struct {
	loglevel uint16
}

func (re View) begin() {
	for k, _ := range prefix_codes {
		prefix_codes[k] += " "
	}
	re.log(LL_DEBUG, "Starting View object.")
}

func (re View) end() {
	re.log(LL_DEBUG, "Stopping View object.")
}

func (re View) output(s string) {
	fmt.Println(s)
	os.Stdout.Sync()
}

func (re View) log(i uint16, s string) {
	if re.loglevel >= i {
		fmt.Fprintln(os.Stderr, prefix_codes[i] + s)
		os.Stderr.Sync()
	}
}
