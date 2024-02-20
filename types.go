package main

const (
	TOKEN_START     = "TOKEN_START"
	TOKEN_STRING    = "TOKEN_STRING"
	TOKEN_NUMBER    = "TOKEN_NUMBER"
	TOKEN_BOOLEAN   = "TOKEN_BOOLEAN"
	TOKEN_COL       = "TOKEN_COL"
	TOKEN_COMMA     = "TOKEN_COMMA"
	TOKEN_OBJ_START = "TOKEN_OBJ_START"
	TOKEN_OBJ_END   = "TOKEN_OBJ_END"
	TOKEN_ARR_START = "TOKEN_ARR_START"
	TOKEN_ARR_END   = "TOKEN_ARR_END"
	TOKEN_NULL      = "TOKEN_NULL"
)

type TokenType string

type Token struct {
	Type  TokenType
	Value string
	Line  int
}

type Lexer struct {
	Tokens []Token
	Buffer []byte
	Length int
	Line   int
	// TokenIndex int
	// ? Children []Lexer
}

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
	Type JsonNodeType
	// Value can be JsonArray, JsonObject or JsonValue
	Value    any
	Children []JsonNode
}

type Parser struct {
	Root JsonNode
}
