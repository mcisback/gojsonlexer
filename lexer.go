package main

import (
	"fmt"
	"log"
)

func isDigit(c byte) bool {
	return c == '0' || c == '1' || c == '2' || c == '3' || c == '4' || c == '5' || c == '6' || c == '7' || c == '8' || c == '9'
}

func incPos(pos *int, inc int) {
	*pos = *pos + inc
}

func (l *Lexer) createToken(Type TokenType, Value string) {
	token := Token{
		Type:  Type,
		Value: Value,
		Line:  l.Line,
	}

	l.Tokens = append(l.Tokens, token)
}

// func (l *Lexer) getToken() Token {
// 	token := l.Tokens[l.TokenIndex]

// 	return token
// }

// func (l *Lexer) shift() Token {
// 	token := l.getToken()

// 	l.TokenIndex += 1

// 	return token
// }

func (l *Lexer) readString(pos *int) string {
	var literal string = ""
	for {
		c := l.Buffer[*pos]

		// TODO: Handle escaping

		if c == '"' {
			break
		}

		literal += string(c)

		incPos(pos, 1)
	}

	return literal
}

func (l *Lexer) readBoolean(pos *int) string {
	var literal string

	if l.Buffer[*pos] == 't' {
		literal = string(l.Buffer[(*pos):(*pos + 4)])

		if literal != "true" {
			log.Fatalln("JSON: Error Wrong Boolean -> ", *pos, literal, string(l.Buffer[(*pos-50):(*pos+50)]))
		}

		fmt.Println("Found Boolean: ", literal)

		incPos(pos, 4)
	} else if l.Buffer[*pos] == 'f' {
		literal = string(l.Buffer[(*pos):(*pos + 5)])

		if literal != "false" {
			log.Fatalln("JSON: Error Wrong Boolean -> ", *pos, literal, string(l.Buffer[(*pos-50):(*pos+50)]))
		}

		fmt.Println("Found Boolean: ", literal)

		incPos(pos, 5)
	} else {
		fmt.Println("JSON: Literal is: ", literal)
		log.Fatalln("JSON: Error Parsing Boolean -> ", *pos, string(l.Buffer[(*pos-50):(*pos+50)]))
	}

	return literal
}

func (l *Lexer) readNumber(pos *int) string {
	var literal string

	for {
		c := l.Buffer[*pos]

		// TODO: Handle Floats and Exponentials
		if isDigit(c) {
			literal += string(c)
		} else {
			*pos = *pos - 1

			break
		}

		incPos(pos, 1)
	}

	return literal
}

func (l *Lexer) readNull(pos *int) string {
	literal := string(l.Buffer[*pos:(*pos + 4)])

	if literal != "null" {
		fmt.Println("JSON: Literal is: ", literal)
		log.Fatalln("JSON: Error Parsing Null Value -> ", *pos, string(l.Buffer[(*pos-50):(*pos+50)]))
	}

	incPos(pos, 4)

	return literal
}

func (l *Lexer) lexer(buffer []byte) {
	l.Buffer = buffer
	l.Length = len(buffer)

	l.createToken(TOKEN_START, "")

	pos := 0
	l.Line = 0

	for {
		if pos >= (l.Length - 1) {
			fmt.Println("JSON: Finished Lexing -> ", pos, l.Length)
			break
		}

		c := l.Buffer[pos]

		fmt.Println("Processing byte: ", string(c))

		switch c {
		// use recursion ? lexer.childern = l.lexer(l.Buffer[i:])
		case '\n', ' ', '\r', '\t':
			if c == '\n' {
				l.Line++
			}

			fmt.Println("Found newline space or tab, skipping...")
			break
		case '{':
			l.createToken(TOKEN_OBJ_START, "{")
		case '[':
			l.createToken(TOKEN_ARR_START, "[")
		case '}':
			l.createToken(TOKEN_OBJ_END, "}")
		case ']':
			l.createToken(TOKEN_ARR_END, "]")
		case ',':
			fmt.Println("Lexer: Found TOKEN_COMMA")

			l.createToken(TOKEN_COMMA, ",")
		case '"':
			pos++

			literal := l.readString(&pos)

			l.createToken(TOKEN_STRING, literal)

			fmt.Println("After readString: ", pos, string(c), literal)
		case ':':
			fmt.Println("Found TOKEN_COL", pos, string(c))

			l.createToken(TOKEN_COL, ":")
		case 't', 'f':
			fmt.Println("Lexer: Found TOKEN_BOOLEAN ?")
			fmt.Println("Found t or f, maybe Boolean: ", string(l.Buffer[pos:(pos+5)]))
			literal := l.readBoolean(&pos)

			pos--

			fmt.Println("After readBoolean: ", pos, literal)

			l.createToken(TOKEN_BOOLEAN, literal)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			literal := l.readNumber(&pos)

			l.createToken(TOKEN_NUMBER, literal)
		case 'n':
			literal := l.readNull(&pos)

			l.createToken(TOKEN_NULL, literal)
		default:
			// fmt.Println("Char, Pos is: ", c, pos)
			log.Fatalln("Error Lexing JSON at: ", l.Line, ":", pos, string(l.Buffer[(pos-10):(pos+100)]))
		}

		pos++
	}
}
