package main

import (
	"fmt"
	"log"
)

const (
	JSON_VALUE_STRING = "JSON_VALUE_STRING"
	JSON_VALUE_NUMBER = "JSON_VALUE_NUMBER"
	JSON_VALUE_BOOLEAN = "JSON_VALUE_BOOLEAN"
	JSON_VALUE_NULL = "JSON_VALUE_NULL"
)

type JsonValueType string

const (
	JSON_NODE_ROOT = "JSON_NODE_ROOT"
	JSON_NODE_VALUE = "JSON_NODE_VALUE"
	JSON_NODE_OBJECT = "JSON_NODE_OBJECT"
	JSON_NODE_ARRAY = "JSON_NODE_ARRAY"
)

type JsonNodeType string

type JsonPointer interface {}

type JsonNode struct {
	Type JsonNodeType
	Value JsonPointer
	Children []JsonNode
}

type JsonObject struct {
	Pairs []JsonPair
}

type JsonArray struct {
	Items []JsonValue
}

type JsonValue struct {
	Type JsonValueType
	Value string
}

type JsonPair struct {
	Key string
	Value JsonValue
}

type Parser struct {
	Root JsonNode
}

func (p *Parser) parseArray(l* Lexer, pos *int) JsonNode {
	
	node := JsonNode{
		Type: JSON_NODE_ARRAY,
	}

	var newValueNode *JsonNode
	
	for {
		if *pos >= len(l.Tokens) {
			log.Fatalln("Error Processing JSON Array")
		}

		token := l.Tokens[*pos]

		fmt.Println("Processing Array Token: ", *pos, token)

		switch(token.Type) {
		case TOKEN_ARR_START:
			fmt.Println("Entering Subarray: ", *pos, token)

			*pos = *pos + 1

			newNode := p.parseArray(l, pos)
			
			node.Children = append(node.Children, newNode) 

			// Così funziona ma non capisco perche'
			// Forse perchè in TOKEN_ARR_END non incrementa pos ?
			*pos = *pos + 1

		// TODO: Object start
		//case TOKEN_OBJ_START:

		case TOKEN_STRING:
			newValueNode = &JsonNode{
				Type: JSON_NODE_VALUE,
				Value: JsonValue{
					Type: JSON_VALUE_STRING,
					Value: token.Value,
				},
			} 

			fmt.Println("Children: ", node.Children)
		case TOKEN_NUMBER:
			newValueNode = &JsonNode{
				Type: JSON_NODE_VALUE,
				Value: JsonValue{
					Type: JSON_VALUE_NUMBER,
					Value: token.Value,
				},
			} 
		case TOKEN_COMMA:
			fmt.Println("Found COMMA !")
			fmt.Println("NewValueNode: ", *newValueNode)
			node.Children = append(node.Children, *newValueNode)
		case TOKEN_ARR_END:
			fmt.Println("Array End: ", node)

			return node
		// case TOKEN_COMMA
		}

		*pos = *pos + 1
	}

	return node
}

func (p *Parser) parse(l *Lexer) {

	var currentNode *JsonNode

	pos := 0

	for {
		if pos >= len(l.Tokens) {
			break
		}

		token := l.Tokens[pos]

		fmt.Println("Lexer Token: ", token.Type, token.Value)

		switch(token.Type) {
		case TOKEN_START:
			currentNode = &p.Root
			currentNode.Type = JSON_NODE_ROOT
		//	p.parseObject(&pos)
		case TOKEN_ARR_START:
			fmt.Println("Parsing Array", pos, token.Type)
			pos++

			newNode := p.parseArray(l, &pos)

			currentNode.Children = append(currentNode.Children, newNode)
		case TOKEN_STRING, TOKEN_BOOLEAN, TOKEN_NUMBER, TOKEN_NULL:
			newNode := JsonNode{
				Type: JSON_NODE_VALUE,
				Value: token.Value,
			}

			currentNode.Children = append(currentNode.Children, newNode)
		}

		pos++
	}
}

func printParser(node *JsonNode, level int) {
	for i := 0; i < level; i++ {
		fmt.Print("-")
	}

	fmt.Println("Node Type: ", node.Type)
	
	for i := 0; i < level; i++ {
		fmt.Print("-")
	}
	
	fmt.Println("Node Value: ", node.Value)

	if len(node.Children) > 0 {
		for _, child := range node.Children {
			printParser(&child, level + 1)
		}
	}
}
