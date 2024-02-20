package main

import (
	"fmt"
	"log"
)

const (
	JSON_VALUE_STRING  = "JSON_VALUE_STRING"
	JSON_VALUE_NUMBER  = "JSON_VALUE_NUMBER"
	JSON_VALUE_BOOLEAN = "JSON_VALUE_BOOLEAN"
	JSON_VALUE_NULL    = "JSON_VALUE_NULL"
)

type JsonValueType string

const (
	JSON_NODE_ROOT   = "JSON_NODE_ROOT"
	JSON_NODE_VALUE  = "JSON_NODE_VALUE"
	JSON_NODE_OBJECT = "JSON_NODE_OBJECT"
	JSON_NODE_ARRAY  = "JSON_NODE_ARRAY"
	JSON_NODE_PAIR   = "JSON_NODE_ARRAY"
)

type JsonNodeType string

type JsonPair map[string]JsonNode

type JsonObject struct {
	Pairs JsonPair
}

type JsonArray struct {
	Items []JsonValue
}

type JsonValue struct {
	Type  JsonValueType
	Value string
}

type JsonNode struct {
	Type     JsonNodeType
	Value    any
	Children []JsonNode
}

type Parser struct {
	Root JsonNode
}

func (p *Parser) parseArray(l *Lexer, pos *int) JsonNode {

	node := JsonNode{
		Type: JSON_NODE_ARRAY,
	}

	// var newValueNode *JsonNode = nil

	for {
		if *pos >= len(l.Tokens) {
			log.Fatalln("Error Processing JSON Array")
		}

		token := l.Tokens[*pos]

		fmt.Println("Processing Array Token: ", *pos, token)

		switch token.Type {
		case TOKEN_ARR_START:
			fmt.Println("Entering Subarray: ", *pos, token)

			*pos = *pos + 1

			newNode := p.parseArray(l, pos)

			node.Children = append(node.Children, newNode)

			// Così funziona ma non capisco perche'
			// Forse perchè in TOKEN_ARR_END non incrementa pos ?
			*pos = *pos + 1

		// TODO: Object start
		case TOKEN_OBJ_START:
			fmt.Println("Entering Subobject inside Array: ", l.Line, ":", *pos, token)

			*pos = *pos + 1

			newNode := p.parseObject(l, pos)

			node.Children = append(node.Children, newNode)

			// Così funziona ma non capisco perche'
			// Forse perchè in TOKEN_ARR_END non incrementa pos ?
			// *pos = *pos + 1

		case TOKEN_STRING, TOKEN_BOOLEAN, TOKEN_NUMBER, TOKEN_NULL:
			nextToken := l.Tokens[(*pos + 1)]

			fmt.Println("Found VALUE: ", token.Type, token.Value)
			fmt.Println("Found VALUE, NEXT TOKEN: ", nextToken.Type, nextToken.Value)

			newNode := p.parseValue(l, pos)
			// newValueNode = &newNode

			fmt.Println("VALUE NewValueNode: ", newNode.Type, newNode.Value)

			node.Children = append(node.Children, newNode)

		// case TOKEN_COMMA:
		// 	if newValueNode == nil {
		// 		fmt.Println("Error TOKEN_COMMA JSON Array at ", l.Line, ":", *pos)

		// 		break
		// 	}

		// 	fmt.Println("Found COMMA !")
		// 	fmt.Println("NewValueNode: ", *newValueNode)
		// 	node.Children = append(node.Children, *newValueNode)

		// 	newValueNode = nil
		// case TOKEN_COL:
		// 	log.Fatalln("Error TOKEN_COL Inside JSON Array at: ", l.Line, ":", *pos)
		case TOKEN_ARR_END:
			fmt.Println("Array End: ", node)

			return node
			// case TOKEN_COMMA
		}

		*pos = *pos + 1
	}

	// return node
}

func (p *Parser) parseObject(l *Lexer, pos *int) JsonNode {
	node := JsonNode{
		Type: JSON_NODE_OBJECT,
	}

	object := &JsonObject{
		Pairs: make(JsonPair),
	}

	// var newValueNode *JsonNode = nil
	// var newPair *JsonPair = nil
	// var objKey string = ""

	for {
		if *pos >= len(l.Tokens) {
			log.Fatalln("Error Processing JSON Object at: ", *pos, ":", l.Line)
		}

		token := l.Tokens[*pos]

		fmt.Println("Processing Object Token: ", *pos, token)

		switch token.Type {
		case TOKEN_OBJ_START:
			log.Fatalln("Error: Nested Object Without Key")

		// TODO: Object start
		//case TOKEN_OBJ_START:

		case TOKEN_STRING: // Prende il valore fino a : e lo mette in key
			// keyNode := p.parseValue(l, pos)
			// newValueNode = &newNode

			// Get Key
			objKey := token.Value

			// Check Column
			*pos = *pos + 1

			token := l.Tokens[*pos]

			if token.Type != TOKEN_COL {
				log.Fatalln("Missing TOKEN_COL after OBJ Key at ", token.Line, ":", *pos)
			}

			// Get Value

			*pos = *pos + 1

			token = l.Tokens[*pos]

			valueNode := &JsonNode{}

			// TODO Process array and object
			// For now works just with values
			switch token.Type {
			case TOKEN_OBJ_START:
				fmt.Println("Entering Child Object: ", *pos, token)

				*pos = *pos + 1

				newNode := p.parseObject(l, pos)

				valueNode = &newNode

				// Così funziona ma non capisco perche'
				// Forse perchè in TOKEN_OBJ_END non incrementa pos ?
				*pos = *pos + 1
			case TOKEN_ARR_START:
				fmt.Println("Entering Subarray in Object: ", *pos, token)

				*pos = *pos + 1

				newNode := p.parseArray(l, pos)

				valueNode = &newNode

				// Così funziona ma non capisco perche'
				// Forse perchè in TOKEN_ARR_END non incrementa pos ?
				*pos = *pos + 1
			case TOKEN_STRING, TOKEN_NULL, TOKEN_NUMBER, TOKEN_BOOLEAN:
				newNode := p.parseValue(l, pos)

				valueNode = &newNode
			}

			fmt.Println("New PAIR: ", objKey, valueNode)

			object.Pairs[objKey] = *valueNode
		// case TOKEN_COL: // Prende il valore fino a COMMA e crea nuova pair

		// case TOKEN_COMMA: // Aggiunge la pair all'oggetto
		// 	if newValueNode == nil {
		// 		log.Fatalln("Error TOKEN_COMMA JSON Array")
		// 	}

		// 	fmt.Println("Found COMMA !")
		// 	fmt.Println("NewValueNode: ", *newValueNode)
		// 	node.Children = append(node.Children, *newValueNode)

		// 	newValueNode = nil
		case TOKEN_OBJ_END:
			node.Value = *object

			fmt.Println("Object End: ", node)

			return node
			// case TOKEN_COMMA
		}

		*pos = *pos + 1
	}
}

func (p *Parser) parseValue(l *Lexer, pos *int) JsonNode {
	var valueType JsonValueType

	token := l.Tokens[*pos]

	switch token.Type {
	case TOKEN_STRING:
		valueType = JSON_VALUE_STRING
	case TOKEN_BOOLEAN:
		valueType = JSON_VALUE_BOOLEAN
	case TOKEN_NUMBER:
		valueType = JSON_VALUE_NUMBER
	case TOKEN_NULL:
		valueType = JSON_VALUE_NULL
	}

	newNode := JsonNode{
		Type: JSON_NODE_VALUE,
		Value: JsonValue{
			Type:  valueType,
			Value: token.Value,
		},
	}

	return newNode
}

func (p *Parser) parse(l *Lexer) {

	var rootNode *JsonNode

	pos := 0

	for {
		if pos >= len(l.Tokens) {
			break
		}

		token := l.Tokens[pos]

		fmt.Println("Lexer Token: ", token.Type, token.Value)

		switch token.Type {
		case TOKEN_START:
			rootNode = &p.Root
			rootNode.Type = JSON_NODE_ROOT
		//	p.parseObject(&pos)
		case TOKEN_ARR_START:
			fmt.Println("Parsing Array", pos, token.Type)
			pos++

			newNode := p.parseArray(l, &pos)

			rootNode.Children = append(rootNode.Children, newNode)
		case TOKEN_OBJ_START:
			fmt.Println("Parsing Array", pos, token.Type)
			pos++

			newNode := p.parseObject(l, &pos)

			rootNode.Children = append(rootNode.Children, newNode)
		case TOKEN_STRING, TOKEN_BOOLEAN, TOKEN_NUMBER, TOKEN_NULL:
			newNode := p.parseValue(l, &pos)

			rootNode.Children = append(rootNode.Children, newNode)
		}

		pos++
	}
}

func printRepeatStr(str string, level int) {
	// fmt.Print(level)

	for i := 0; i < level; i++ {
		fmt.Print(str)
	}
}

func printParser(node *JsonNode, level int) {
	printRepeatStr("-", level)

	// TODO Print Pairs

	if node.Type == JSON_NODE_ARRAY {
		fmt.Println("Array ", "(", len(node.Children), ") ", "([")
	} else if node.Type == JSON_NODE_OBJECT {
		fmt.Println("Object ", "(", len(node.Children), ") ", "{")
	} else {
		fmt.Println("Node(t:", node.Type, ", v:", node.Value, ", c:", len(node.Children), ")")
	}

	if len(node.Children) > 0 {
		if node.Type == JSON_NODE_ARRAY {
			for _, child := range node.Children {
				printParser(&child, level+1)
			}
		}
	}

	if node.Type == JSON_NODE_OBJECT {
		object := node.Value.(JsonObject)

		for key := range object.Pairs {
			pair := object.Pairs[key]

			printParser(&pair, level+1)
		}
	}

	if node.Type == JSON_NODE_ARRAY {
		printRepeatStr("-", level)

		fmt.Println("])")
	} else if node.Type == JSON_NODE_OBJECT {
		printRepeatStr("-", level)

		fmt.Println("}")
	}
}

func printJson(node *JsonNode, level int, key string) {
	separator := "\t"

	printRepeatStr(separator, level)

	// TODO Print Pairs

	if node.Type == JSON_NODE_ARRAY {
		fmt.Println("[")
	} else if node.Type == JSON_NODE_OBJECT {
		fmt.Println("{")
	} else {
		if key != "" {
			fmt.Print(key, ": ")
		}

		switch node.Value.(type) {
		case JsonValue:
			value := node.Value.(JsonValue)

			switch value.Type {
			case JSON_VALUE_BOOLEAN, JSON_VALUE_NUMBER, JSON_VALUE_NULL:
				fmt.Print(value.Value, ",")
			case JSON_VALUE_STRING:
				fmt.Print("\"", value.Value, "\",")
			}

		}

		fmt.Println("")
	}

	if len(node.Children) > 0 {
		if node.Type == JSON_NODE_ARRAY {
			for _, child := range node.Children {
				printJson(&child, level+1, "")
			}
		}
	}

	if node.Type == JSON_NODE_OBJECT {
		object := node.Value.(JsonObject)

		for key := range object.Pairs {
			pair := object.Pairs[key]

			printJson(&pair, level+1, key)
		}
	}

	if node.Type == JSON_NODE_ARRAY {
		printRepeatStr(separator, level)

		fmt.Println("]")
	} else if node.Type == JSON_NODE_OBJECT {
		printRepeatStr(separator, level)

		fmt.Println("}")
	}
}
