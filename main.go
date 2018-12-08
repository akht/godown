package main

import (
	"godown/converter"
	"os"
	"os/user"
)

func main() {
	_, err := user.Current()
	if err != nil {
		panic(err)
	}

	converter.Convert(os.Stdin, os.Stdout)
}
