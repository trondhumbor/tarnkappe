package main

import (
	"fmt"
	"os"

	"github.com/trondhumbor/tarnkappe/internal/tarnkappe"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("hide inpath outpath hidepath")
	}

	result, err := tarnkappe.Hide(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
