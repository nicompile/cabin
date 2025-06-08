package main

import (
	"log"
	"os"

	"github.com/nicompile/cabin/internal/build"
)

func main() {
	url := os.Args[1]

	err := build.Build(url, "serverless-api")
	if err != nil {
		log.Fatal(err)
	}
}
