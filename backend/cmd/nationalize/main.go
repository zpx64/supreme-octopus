package main

import (
	"fmt"
	"os"

	"github.com/k0kubun/pp"
	"github.com/ssleert/tzproj/pkg/nationalize"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no arg")
		os.Exit(2)
	}
	c, err := nationalize.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := c.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	pp.Println(out)
}
