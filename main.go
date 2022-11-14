package main

import (
	"fmt"

	"github.com/xilu0/transctl/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Print(err.Error())
	}
}
