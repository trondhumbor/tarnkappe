package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/trondhumbor/tarnkappe/internal/tarnkappe"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("reveal inpath outpath length")
	}

	l, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(err)
	}

	tarnkappe.Reveal(os.Args[1], os.Args[2], l)
}
