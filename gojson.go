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

func (json *GoJson) get(path string, rootNode *JsonNode) *JsonNode {

	var node *JsonNode = rootNode
	var arrayValueRegex = regexp.MustCompile(`^\[(\d+)]`)
	var arrayByIndexRegex = regexp.MustCompile(`^(\d+)`)
	var objectValueRegex = regexp.MustCompile(`^\[\"([^"]+)\"]`)
	var keyRegex = regexp.MustCompile(`^(\w+)`)
	// var lenRegex = regexp.MustCompile(`^len\(([^\)]+)\)`)
	// var getKeysRegex = regexp.MustCompile(`^keys\(([^\)]+)\)`)

	if node == nil {
		node = &json.parser.Root
	}

	funcs := strings.Split(path, "|")

	// fmt.Println("funcs: ", funcs)

	for i := 0; i < len(funcs); i++ {
		f := funcs[i]

		// fmt.Println("funcs: ", f)

		if f == "len" {

			node = json.get(funcs[i+1], nil)

			// fmt.Println("Matched len()")

			var length int = 0

			switch node.Type {
			case JSON_NODE_ARRAY:
				length = len(node.Children)

			case JSON_NODE_OBJECT:
				object := node.Value.(JsonObject)

				length = len(object.Pairs)
			default:
				log.Fatalln("Error, matching for length in a JsonValue")
			}

			value := JsonNode{
				Type: JSON_NODE_VALUE,
				Value: JsonValue{
					Type:  JSON_VALUE_NUMBER,
					Value: fmt.Sprintf("%d", length),
				},
			}

			return &value
		}

		if f == "keys" {

			node = json.get(funcs[i+1], nil)

			// fmt.Println("Matched keys()")

			switch node.Type {
			case JSON_NODE_OBJECT:
				object := node.Value.(JsonObject)

				tempNode := JsonNode{
					Type: JSON_NODE_ARRAY,
				}

				// fmt.Println("[")

				for key := range object.Pairs {
					value := JsonNode{
						Type: JSON_NODE_VALUE,
						Value: JsonValue{
							Type:  JSON_VALUE_STRING,
							Value: key,
						},
					}

					tempNode.Children = append(tempNode.Children, value)
					// fmt.Print("\t\"", key, "\",\n")
				}

				printJson(&tempNode, 0, "")

				// fmt.Println("]")
			default:
				log.Fatalln("Error, getting keys in a non JsonObject")
			}

			// what if keys($.[0]).[0] ?
			return node

		}
	}

	parts := strings.Split(path, ".")

	for i, part := range parts {
		// fmt.Println("PART: ", part)

		// if node != nil {
		// 	fmt.Println("CURRENT_NODE: ", *node)
		// }

		if part == "$" && i == 0 {
			if len(node.Children) <= 0 {
				log.Fatalln("Error, Array out of range")
			}

			node = &node.Children[0]

			// fmt.Println("QUERY_ROOT_NODE: ", part, *node)
		} else if part == "[]" {
			fmt.Println("node: ", node)
			fmt.Println("node.Children: ", node.Children)

			if len(node.Children) <= 0 {
				log.Fatalln("JSON Error, Array out of range")
			}

			resultArray := JsonNode{
				Type: JSON_NODE_ARRAY,
			}

			for _, n := range node.Children {
				joined := strings.Join(parts[(i+1):], ".")

				if joined == "" {
					return node
				}

				// fmt.Println("[")

				// value := JsonNode{
				// 	Type: JSON_NODE_VALUE,
				// 	Value: JsonValue{
				// 		Type:  JSON_VALUE_STRING,
				// 		Value: key,
				// },

				fmt.Println("[] Joined Parts: ", joined)

				value := json.get(joined, &n)

				resultArray.Children = append(resultArray.Children, *value)
				// fmt.Print("\t\"", key, "\",\n")

				fmt.Println("ResultArray: ", resultArray)

			}

			return &resultArray

		} else if match := arrayByIndexRegex.FindStringSubmatch(part); len(match) > 1 {
			if node.Type != JSON_NODE_ARRAY {
				log.Fatalln("Error, matching for an index in non JSON Array")
			}

			index, _ := strconv.ParseInt(match[1], 10, 32)

			if int(index) >= len(node.Children) {
				log.Fatalln("Error, Array out of range")
			}

			node = &node.Children[index]

			// fmt.Println("QUERY_ARRAY_VALUE: ", match[1], index)
		} else if match := arrayValueRegex.FindStringSubmatch(part); len(match) > 1 {
			if node.Type != JSON_NODE_ARRAY {
				log.Fatalln("Error, matching for an index in non JSON Array")
			}

			index, _ := strconv.ParseInt(match[1], 10, 32)

			if int(index) >= len(node.Children) {
				log.Fatalln("Error, Array out of range")
			}

			node = &node.Children[index]

			// fmt.Println("QUERY_ARRAY_VALUE: ", match[1], index)
		} else if match := objectValueRegex.FindStringSubmatch(part); len(match) > 1 {

			if node.Type != JSON_NODE_OBJECT {
				log.Fatalln("Error, matching for a key in non JSON Object")
			}

			key := match[1]

			// object := node.Value.(JsonObject)
			object := node.Value.(JsonObject)
			value, ok := object.Pairs[key]

			if !ok {
				log.Fatalln("Error, key: \"", key, "\" not found in JSON Object")
			}

			node = &value

			//node = object.Pairs[index]

			// fmt.Println("QUERY_OBJECT_VALUE: ", match[1], key, value)
		} else if match := keyRegex.FindStringSubmatch(part); len(match) > 1 {

			if node.Type != JSON_NODE_OBJECT {
				log.Fatalln("Error, matching for a key in non JSON object")
			}

			key := match[1]

			// object := node.Value.(JsonObject)
			object := node.Value.(JsonObject)
			value, ok := object.Pairs[key]

			if !ok {
				log.Fatalln("Error, key: \"", key, "\" not found in JSON Object")
			}

			node = &value

			//node = object.Pairs[index]

			// fmt.Println("QUERY_OBJECT_KEY: ", match[1], key, value)
		} else {
			log.Fatalln("Error, cannot undestrand: \"", part, "\"")

			// return node
		}

		if node == nil {
			log.Fatalln("Missing $ ?")
		}
	}

	return node

}

func (json *GoJson) outputJson(node *JsonNode) {
	if node.Type == JSON_NODE_ARRAY || node.Type == JSON_NODE_OBJECT {
		printJson(node, 0, "")
	} else if node.Type == JSON_NODE_VALUE {
		printJsonValue(node, true)
	}
}
