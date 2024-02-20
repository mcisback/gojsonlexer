package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type GoJson struct {
	// rawData []bytes
	lexer  *Lexer
	parser *Parser
}

func (json *GoJson) fromFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	json.lexer = &Lexer{}
	json.parser = &Parser{}

	json.lexer.lexer(bytes)

	// fmt.Println("Lexer: ", lexer.Tokens)
	fmt.Println("Starting Parser: ")

	json.parser.parse(json.lexer)

	fmt.Println("Parser: ", *json.parser)

	// fmt.Println("Printin Parser: ")

	// printParser(&json.parser.Root, 0)
}

func (json *GoJson) get(path string) {

	var node *JsonNode = nil
	var arrayValueRegex = regexp.MustCompile(`^\[(\d+)]`)
	var objectValueRegex = regexp.MustCompile(`^\[\"([^"]+)\"]`)
	var keyRegex = regexp.MustCompile(`^(\w+)`)

	parts := strings.Split(path, ".")

	for _, part := range parts {
		// fmt.Println("PART: ", part)

		// if node != nil {
		// 	fmt.Println("CURRENT_NODE: ", *node)
		// }

		if part == "$" {
			node = &json.parser.Root

			// fmt.Println("QUERY_ROOT_NODE: ", part, *node)
		} else if part == "[]" {
			node = &node.Children[0]
			// fmt.Println("QUERY_ARRAY: ", part)
		} else if match := arrayValueRegex.FindStringSubmatch(part); len(match) > 1 {
			index, _ := strconv.ParseInt(match[1], 10, 32)

			if int(index) >= len(node.Children) {
				log.Fatalln("Array out of range")
			}

			node = &node.Children[index]

			// fmt.Println("QUERY_ARRAY_VALUE: ", match[1], index)
		} else if match := objectValueRegex.FindStringSubmatch(part); len(match) > 1 {
			key := match[1]

			// object := node.Value.(JsonObject)
			object := node.Value.(JsonObject)
			value := object.Pairs[key]

			node = &value

			//node = object.Pairs[index]

			// fmt.Println("QUERY_OBJECT_VALUE: ", match[1], key, value)
		} else if match := keyRegex.FindStringSubmatch(part); len(match) > 1 {
			key := match[1]

			// object := node.Value.(JsonObject)
			object := node.Value.(JsonObject)
			value := object.Pairs[key]

			node = &value

			//node = object.Pairs[index]

			// fmt.Println("QUERY_OBJECT_KEY: ", match[1], key, value)
		}

		if node == nil {
			log.Fatalln("Missing $ ?")
		}
	}

	if node.Type == JSON_NODE_ARRAY || node.Type == JSON_NODE_OBJECT {
		printJson(node, 0, "")
	} else if node.Type == JSON_NODE_VALUE {
		printJsonValue(node, "", true)
	}

}
