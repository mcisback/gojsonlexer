package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Missing json query")
	}

	if len(os.Args) <= 2 {
		log.Fatal("Missing json file")
	}

	var jsonPath = os.Args[1]
	var jsonFilename = os.Args[2]

	json := GoJson{}

	json.fromFile(jsonFilename)

	json.get(jsonPath)

}
