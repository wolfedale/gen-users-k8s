package main

import (
	"log"
	"os"

	"./cmd"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Wrong number of arguments.")
		cmd.NewRootCmd().Help()
		os.Exit(2)
	}
	cmd.Run()
}
