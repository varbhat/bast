// Lexer for BAST
package main

import (
	"fmt"
	"os"
	"unicode"
)

// Lexer Struct
type Lexer struct {
	source  string
	curChar string
	curPos  int
	length  int
}

// Initialize the lexer
func (lexer *Lexer) initlexer(input string) {
	lexer.source = input + "\n"
	lexer.length = len(lexer.source)
	lexer.curChar = ""
	lexer.curPos = -1
	lexer.nextChar()
}

// Process next Character
func (lexer *Lexer) nextChar() {
	lexer.curPos += 1
	if lexer.curPos >= lexer.length {
		lexer.curChar = string("\u0000") // EOF
	} else {
		lexer.curChar = string(lexer.source[lexer.curPos])
	}
}

// Return the lookahead Character
func (lexer *Lexer) peek() string {
	if lexer.curPos+1 >= len(lexer.source) {
		return string("\u0000")
	} else {
		return string(lexer.source[lexer.curPos+1])
	}
}

// Skip Whitespace except newlines , which we will use to indicate the end of a statement.
func (lexer *Lexer) skipWhitespace() {
	if lexer.curChar == "\n" {
		return
	}
	for unicode.IsSpace(rune(lexer.curChar[0])) == true {
		lexer.nextChar()
	}
}

// Skip the New Line beginning with #
func (lexer *Lexer) skipComments() {
	if lexer.curChar == "#" {
		for lexer.curChar != "\n" {
			lexer.nextChar()
		}
	}
}

// Tokens

// Token Type Struct
type TokenT struct {
	idname string
	idno   int
}

// Function To get Preinitialized TokenT with type of idname returned only for keywords
func mul_initTokenT(varidname string) (TokenT, string) {
	switch varidname {
	case "EOF":
		return TokenT{idname: varidname, idno: -1}, ""
	case "NEWLINE":
		return TokenT{idname: varidname, idno: 0}, ""
	case "NUMBER":
		return TokenT{idname: varidname, idno: 1}, ""
	case "IDENT":
		return TokenT{idname: varidname, idno: 2}, ""
	case "STRING":
		return TokenT{idname: varidname, idno: 3}, ""
	// Keywords
	case "LABEL":
		return TokenT{idname: varidname, idno: 101}, "LABEL"
	case "GOTO":
		return TokenT{idname: varidname, idno: 102}, "GOTO"
	case "PRINT":
		return TokenT{idname: varidname, idno: 103}, "PRINT"
	case "INPUT":
		return TokenT{idname: varidname, idno: 104}, "INPUT"
	case "LET":
		return TokenT{idname: varidname, idno: 105}, "LET"
	case "IF":
		return TokenT{idname: varidname, idno: 106}, "IF"
	case "THEN":
		return TokenT{idname: varidname, idno: 107}, "THEN"
	case "ENDIF":
		return TokenT{idname: varidname, idno: 108}, "ENDIF"
	case "WHILE":
		return TokenT{idname: varidname, idno: 109}, "WHILE"
	case "REPEAT":
		return TokenT{idname: varidname, idno: 110}, "REPEAT"
	case "ENDWHILE":
		return TokenT{idname: varidname, idno: 111}, "ENDWHILE"
	// Operators
	case "EQ":
		return TokenT{idname: varidname, idno: 201}, ""
	case "PLUS":
		return TokenT{idname: varidname, idno: 202}, ""
	case "MINUS":
		return TokenT{idname: varidname, idno: 203}, ""
	case "ASTERISK":
		return TokenT{idname: varidname, idno: 204}, ""
	case "SLASH":
		return TokenT{idname: varidname, idno: 205}, ""
	case "EQEQ":
		return TokenT{idname: varidname, idno: 206}, ""
	case "NOTEQ":
		return TokenT{idname: varidname, idno: 207}, ""
	case "LT":
		return TokenT{idname: varidname, idno: 208}, ""
	case "LTEQ":
		return TokenT{idname: varidname, idno: 210}, ""
	case "GT":
		return TokenT{idname: varidname, idno: 211}, ""
	case "GTEQ":
		return TokenT{idname: varidname, idno: -1}, ""
	default:
		return TokenT{idname: "", idno: -10000}, ""
	}
}

// We don't need second value of mul_initTokenT always
func initTokenT(varidname string) TokenT {
	tktobj, _ := mul_initTokenT(varidname)
	return tktobj
}

// Struct for Token
type Token struct {
	text string
	kind TokenT
}

func initToken(curchar string, typevar TokenT) Token {
	return Token{
		text: curchar,
		kind: typevar,
	}
}

func (lexer *Lexer) getToken() Token {
	lexer.skipWhitespace()
	lexer.skipComments()
	var tokenvar Token

	// Check the first character of this token to see if we can decide what it is.
	// If it is a multiple character operator (e.g.,!=),number,identifier,or keyword then we will proceed the rest.
	if lexer.curChar == "+" {
		tokenvar = initToken("+", initTokenT("PLUS"))
	} else if lexer.curChar == "-" {
		tokenvar = initToken("-", initTokenT("MINUS"))
	} else if lexer.curChar == "*" {
		tokenvar = initToken("*", initTokenT("ASTERISK"))
	} else if lexer.curChar == "/" {
		tokenvar = initToken("/", initTokenT("SLASH"))
	} else if lexer.curChar == "\n" {
		tokenvar = initToken("\n", initTokenT("NEWLINE"))
	} else if lexer.curChar == "\u0000" {
		tokenvar = initToken("\u0000", initTokenT("EOF"))
	} else if lexer.curChar == "=" {
		// Check whether this token is = or ==
		if lexer.peek() == "=" {
			lexer.nextChar()
			tokenvar = initToken("==", initTokenT("EQEQ"))
		} else {
			tokenvar = initToken("=", initTokenT("EQ"))
		}
	} else if lexer.curChar == ">" {
		// Check whether this token is > or >=
		if lexer.peek() == "=" {
			lexer.nextChar()
			tokenvar = initToken(">=", initTokenT("GTEQ"))
		} else {
			tokenvar = initToken(">", initTokenT("GT"))
		}
	} else if lexer.curChar == "<" {
		// Check whether this token is < or <=
		if lexer.peek() == "=" {
			lexer.nextChar()
			tokenvar = initToken("<=", initTokenT("LTEQ"))
		} else {
			tokenvar = initToken("<", initTokenT("LT"))
		}
	} else if lexer.curChar == "!" {
		// Check whether this token is != or !=
		if lexer.peek() == "=" {
			lexer.nextChar()
			tokenvar = initToken("!=", initTokenT("NOTEQ"))
		} else {
			fmt.Println(term_col_err, "Expected !=, got !", term_col_res)
			os.Exit(1)
		}
	} else if lexer.curChar == "\"" {
		// Get characters between quotations
		lexer.nextChar()
		startPos := lexer.curPos

		for lexer.curChar != "\"" {
			// Don't allow special characters in the string.No escape characters,newlines,tabs,or %
			// We will be using C's printf on this string
			if (lexer.curChar == "\n") || (lexer.curChar == "\r" || (lexer.curChar == "\t") || (lexer.curChar == "\\") || (lexer.curChar == "%")) {
				fmt.Println(term_col_err, "Illegal character in string.", term_col_res)
				os.Exit(1)
			}
			lexer.nextChar()
		}

		tokText := lexer.source[startPos:lexer.curPos] // Get the substring
		tokenvar = initToken(tokText, initTokenT("STRING"))
	} else if unicode.IsDigit(rune(lexer.curChar[0])) == true {
		//Leading Character is a digit, so this must be a number.
		// Get all consecutive digits and decimal if there is one.
		startPos := lexer.curPos
		for (unicode.IsDigit(rune(lexer.peek()[0]))) == true {
			lexer.nextChar()
		}
		if (lexer.peek() == ".") == true {
			lexer.nextChar()

			// Must have atleast one digit after decimal.
			if unicode.IsDigit(rune(lexer.peek()[0])) != true {
				// Error
				fmt.Println(term_col_err, "Illegal Character in number", term_col_res)
				os.Exit(1)
			}
			for unicode.IsDigit(rune(lexer.peek()[0])) == true {
				lexer.nextChar()
			}
		}
		tokText := lexer.source[startPos : lexer.curPos+1]
		tokenvar = initToken(tokText, initTokenT("NUMBER"))
	} else if (unicode.IsLetter(rune(lexer.curChar[0]))) == true {
		//Leading Character is a letter, so this must be a identifier.
		// Get all consecutive alpha numeric characters.
		startPos := lexer.curPos
		for (unicode.IsLetter(rune(lexer.peek()[0])) || (unicode.IsDigit(rune(lexer.peek()[0])))) == true {
			lexer.nextChar()
		}
		tokText := lexer.source[startPos : lexer.curPos+1]
		_, checkkeyword := mul_initTokenT(tokText)
		if checkkeyword == "" {
			tokenvar = initToken(tokText, initTokenT("IDENT"))
		} else {
			tokenvar = initToken(tokText, initTokenT(checkkeyword))
		}

	} else {
		fmt.Println(term_col_err, "Lexing error.Unknown Token:", term_col_res, lexer.curChar)
	}
	lexer.nextChar()
	return tokenvar
}
