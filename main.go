package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
)

func readFile(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename) // just pass the file name
	if err != nil {
   	log.Fatalln(err)
	}

	return bytes
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Missing json filename")
	}

	var jsonFilename = os.Args[1]

	bytes := readFile(jsonFilename)

	lexer := Lexer{}

	lexer.lexer(bytes)

	// fmt.Println("Lexer: ", lexer.Tokens)
	fmt.Println("Starting Parser: ")

	parser := Parser{}

	parser.parse(&lexer)

	fmt.Println("Parser: ", parser)

	printParser(&parser.Root, 0)
}
