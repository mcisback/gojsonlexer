package main

import (
	"log"
)

func (p *Parser) parseArray(l *Lexer, pos *int) JsonNode {

	node := JsonNode{
		Type: JSON_NODE_ARRAY,
	}

	// var newValueNode *JsonNode = nil

	for {
		if *pos >= len(l.Tokens) {
			// fmt.Println("Error Processing JSON Array at: ", *pos, l.Line)

			return node
		}

		token := l.Tokens[*pos]

		// fmt.Println("Processing Array Token: ", *pos, token)

		switch token.Type {
		case TOKEN_ARR_START:
			// fmt.Println("Entering Subarray: ", *pos, token)

			*pos = *pos + 1

			newNode := p.parseArray(l, pos)

			node.Children = append(node.Children, newNode)

			// Così funziona ma non capisco perche'
			// Forse perchè in TOKEN_ARR_END non incrementa pos ?
			*pos = *pos + 1

		case TOKEN_OBJ_START:
			// fmt.Println("Entering Subobject inside Array: ", l.Line, ":", *pos, token)

			*pos = *pos + 1

			newNode := p.parseObject(l, pos)

			node.Children = append(node.Children, newNode)

			// Così funziona ma non capisco perche'
			// Forse perchè in TOKEN_ARR_END non incrementa pos ?
			// *pos = *pos + 1

		case TOKEN_STRING, TOKEN_BOOLEAN, TOKEN_NUMBER, TOKEN_NULL:
			// nextToken := l.Tokens[(*pos + 1)]

			// fmt.Println("Found VALUE: ", token.Type, token.Value)
			// fmt.Println("Found VALUE, NEXT TOKEN: ", nextToken.Type, nextToken.Value)

			newNode := p.parseValue(l, pos)
			// newValueNode = &newNode

			// fmt.Println("VALUE NewValueNode: ", newNode.Type, newNode.Value)

			node.Children = append(node.Children, newNode)

		case TOKEN_ARR_END:
			// fmt.Println("Array End: ", node)

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

		// fmt.Println("Processing Object Token: ", *pos, token)

		switch token.Type {
		case TOKEN_OBJ_START:
			log.Fatalln("Error: Nested Object Without Key")

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

			switch token.Type {
			case TOKEN_OBJ_START:
				// fmt.Println("Entering Child Object: ", *pos, token)

				*pos = *pos + 1

				newNode := p.parseObject(l, pos)

				valueNode = &newNode

				// Così funziona ma non capisco perche'
				// Forse perchè in TOKEN_OBJ_END non incrementa pos ?
				*pos = *pos + 1
			case TOKEN_ARR_START:
				// fmt.Println("Entering Subarray in Object: ", *pos, token)

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

			// fmt.Println("New PAIR: ", objKey, valueNode)

			object.Pairs[objKey] = *valueNode
		case TOKEN_OBJ_END:
			node.Value = *object

			// fmt.Println("Object End: ", node)

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

		// fmt.Println("Lexer Token: ", token.Type, token.Value)

		switch token.Type {
		case TOKEN_START:
			rootNode = &p.Root
			rootNode.Type = JSON_NODE_ROOT
		//	p.parseObject(&pos)
		case TOKEN_ARR_START:
			// fmt.Println("Parsing Array", pos, token.Type)
			pos++

			newNode := p.parseArray(l, &pos)

			rootNode.Children = append(rootNode.Children, newNode)
		case TOKEN_OBJ_START:
			// fmt.Println("Parsing Array", pos, token.Type)
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
