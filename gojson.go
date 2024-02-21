package main

import (
	"fmt"
	"log"
	"maps"
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

func (json *GoJson) fromString(jsonData string) {
	json.lexer = &Lexer{}
	json.parser = &Parser{}

	json.lexer.lexer([]byte(jsonData))

	// fmt.Println("Lexer: ", lexer.Tokens)
	// fmt.Println("Starting Parser: ")

	json.parser.parse(json.lexer)

	// fmt.Println("Parser: ", *json.parser)

	// fmt.Println("Printin Parser: ")

	// printParser(&json.parser.Root, 0)
}

func (json *GoJson) fromFile(filename string) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	json.fromString(string(bytes))
}

// func (json *GoJson) getNodeLength(funcs []string, i int) *JsonNode {
// 	node := json.getNode(funcs[i+1], nil)

// 	// fmt.Println("Matched len()")

// 	var length int = 0

// 	switch node.Type {
// 	case JSON_NODE_ARRAY:
// 		length = len(node.Children)

// 	case JSON_NODE_OBJECT:
// 		object := node.Value.(JsonObject)

// 		length = len(object.Pairs)
// 	default:
// 		log.Fatalln("Error, matching for length in a JsonValue")
// 	}

// 	value := JsonNode{
// 		Type: JSON_NODE_VALUE,
// 		Value: JsonValue{
// 			Type:  JSON_VALUE_NUMBER,
// 			Value: fmt.Sprintf("%d", length),
// 		},
// 	}

// 	return &value
// }

func SplitAny(s string, seps string) []string {
	splitter := func(r rune) bool {
		return strings.ContainsRune(seps, r)
	}
	return strings.FieldsFunc(s, splitter)
}

func (json *GoJson) getNode(path string, rootNode *JsonNode) *JsonNode {

	// fmt.Println("Called getNode() -> ", path)

	var node *JsonNode = rootNode
	var arrayValueRegex = regexp.MustCompile(`^\[(\d+)]`)
	var arrayByIndexRegex = regexp.MustCompile(`^(\d+)`)
	var objectValueRegex = regexp.MustCompile(`^\[\"([^"]+)\"]`)
	var keyRegex = regexp.MustCompile(`^(\w+)`)
	// var lenRegex = regexp.MustCompile(`^len\(([^\)]+)\)`)
	// var getKeysRegex = regexp.MustCompile(`^keys\(([^\)]+)\)`)
	// var exprsRegex = regexp.MustCompile("[\\:\\,\\.\\s]+")
	var isObjectRegex = regexp.MustCompile(`^\{([^\}]+)\}`)

	if node == nil {
		node = &json.parser.Root
	}

	// $.1 + $.2 or values|$.1+values|$.2 or keys|$.1+keys|$.2
	// exprs := SplitAny(path, "+")

	// if len(exprs) > 1 {

	// 	for i := 0; i < len(exprs); i++ {
	// 		expr := exprs[i]

	// 		fmt.Println("EXPR: ", expr)

	// 	}

	// }

	funcs := strings.Split(path, "|")

	// fmt.Println("funcs: ", funcs)

	for i := 0; i < len(funcs); i++ {
		f := funcs[i]

		// fmt.Println("funcs: ", f)

		if f == "len" {
			node = json.getNode(strings.Join(funcs[i+1:], "|"), node)

			fmt.Println("Matched len()", node)

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
		} else if strings.HasPrefix(f, "keys") {

			// fmt.Println("f == keys")
			var index int = -1
			args := strings.Split(f, ".")

			// fmt.Println("keys: ", args)

			if len(args) > 1 {
				i, _ := strconv.ParseInt(args[1], 10, 32)
				index = int(i)
			}

			// fmt.Println("index: ", index)

			node = json.getNode(strings.Join(funcs[i+1:], "|"), nil)

			// fmt.Println("Matched keys()")

			resultArray := JsonNode{
				Type: JSON_NODE_ARRAY,
			}

			switch node.Type {
			case JSON_NODE_OBJECT:
				object := node.Value.(JsonObject)

				// fmt.Println("[")

				for key := range object.Pairs {
					value := JsonNode{
						Type: JSON_NODE_VALUE,
						Value: JsonValue{
							Type:  JSON_VALUE_STRING,
							Value: key,
						},
					}

					resultArray.Children = append(resultArray.Children, value)
					// fmt.Print("\t\"", key, "\",\n")
				}

				// printJson(&resultArray, 0, "")

				// fmt.Println("]")
			default:
				log.Fatalln("Error, getting keys in a non JsonObject")
			}

			// what if keys($.[0]).[0] ?
			if index >= 0 {
				if index >= len(resultArray.Children) {
					log.Fatalln("Keys Array out of range")
				}

				return &resultArray.Children[index]
			}

			return &resultArray
		} else if strings.HasPrefix(f, "values") {

			// fmt.Println("f == keys")
			var index int = -1
			args := strings.Split(f, ".")

			// fmt.Println("keys: ", args)

			if len(args) > 1 {
				i, _ := strconv.ParseInt(args[1], 10, 32)
				index = int(i)
			}

			// fmt.Println("index: ", index)

			node = json.getNode(strings.Join(funcs[i+1:], "|"), nil)

			// fmt.Println("Matched keys()")

			resultArray := JsonNode{
				Type: JSON_NODE_ARRAY,
			}

			switch node.Type {
			case JSON_NODE_OBJECT:
				object := node.Value.(JsonObject)

				// fmt.Println("[")

				for key := range object.Pairs {

					// value := JsonNode{
					// 	Type:  JSON_NODE_VALUE,
					// 	Value: object.Pairs[key],
					// }

					// fmt.Println("key: ", key, object.Pairs[key])

					resultArray.Children = append(resultArray.Children, object.Pairs[key])
					// fmt.Print("\t\"", key, "\",\n")
				}

				// printJson(&resultArray, 0, "")

				// fmt.Println("]")
			default:
				log.Fatalln("Error, getting values in a non JsonObject")
			}

			// what if keys($.[0]).[0] ?
			if index >= 0 {
				if index >= len(resultArray.Children) {
					log.Fatalln("Keys Array out of range")
				}

				// Works badly i dunno why (order is changes when index changes and i dunno why)
				return &resultArray.Children[index]
			}

			return &resultArray
		} else if match := isObjectRegex.FindStringSubmatch(f); len(match) > 1 {
			// Create new json from string and maybe add it to existing json
			// fmt.Println("isObjectRegex match: ", match)

			jsonData := match[0]

			// fmt.Println("isObjectRegex match: ", jsonData)

			node = json.getNode(strings.Join(funcs[i+1:], "|"), nil)

			// TODO: Add append to array, to object, etc...

			j := GoJson{}

			// FIXME: ? Doesn't work without a newline
			j.fromString(jsonData + "\n")

			// fmt.Println()
			// fmt.Println()

			newNode := j.getNode("$", nil)

			object := node.Value.(JsonObject)

			newPairs := newNode.Value.(JsonObject).Pairs

			maps.Copy(object.Pairs, newPairs)

			// fmt.Println(newPairs)

			// fmt.Println(newNode)

			// fmt.Println(object.Pairs)

			return node
		}
		// } else {
		// 	return json.getNode(strings.Join(funcs[i+1:], "|"), node)
		// }
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
			// fmt.Println("node: ", node)
			// fmt.Println("node.Children: ", node.Children)

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

				// fmt.Println("[] Joined Parts: ", joined)

				value := json.getNode(joined, &n)

				resultArray.Children = append(resultArray.Children, *value)
				// fmt.Print("\t\"", key, "\",\n")

				// fmt.Println("ResultArray: ", resultArray)

			}

			return &resultArray

		} else if match := arrayByIndexRegex.FindStringSubmatch(part); len(match) > 1 {
			// fmt.Println("arraybyIndexRegex: ", match[1])

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

	// fmt.Println("Returning node: ", node)

	return node

}

// func (json *GoJson) getNodeValue(node *JsonNode) string {
// 	// if node.Type == JSON_NODE_ARRAY || node.Type == JSON_NODE_OBJECT {
// 	// 	return node
// 	// } else if node.Type == JSON_NODE_VALUE {
// 	// 	value := node.Value.(JsonValue)

// 	// 	return value.Value
// 	// }

// 	value := node.Value.(JsonValue)

// 	return value.Value
// }

func (json *GoJson) outputJson(node *JsonNode) {
	if node.Type == JSON_NODE_ARRAY || node.Type == JSON_NODE_OBJECT {
		printJson(node, 0, "")
	} else if node.Type == JSON_NODE_VALUE {
		printJsonValue(node, true)
	}
}
