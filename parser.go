// Parser for BAST ( Also emits C code)
package main

import (
	"fmt"
	"os"
)

// Function to check element in string slice
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Create Parser Struct
type Parser struct {
	lexer          Lexer
	curToken       Token
	peekToken      Token
	symbols        []string
	labelsDeclared []string
	labelsGotoed   []string
}

// Initialize Parser
func (parser *Parser) initParser(lexer Lexer) {
	parser.lexer = lexer
	parser.curToken = Token{text: "", kind: TokenT{idname: "", idno: -1000}}
	parser.peekToken = Token{text: "", kind: TokenT{idname: "", idno: -1000}}
	parser.symbols = make([]string, 10)
	parser.labelsDeclared = make([]string, 10)
	parser.labelsGotoed = make([]string, 10)
	parser.nextToken()
	parser.nextToken()
}

// Return true if the current token matches
func (parser *Parser) checkToken(kind string) bool {
	return (kind == parser.curToken.kind.idname)
}

// Return true if the next token matches
func (parser *Parser) checkPeek(kind string) bool {
	return (kind == parser.peekToken.kind.idname)
}

// Try to match current token.If not,error.Advance the current token
func (parser *Parser) matchToken(kind string) {
	if parser.checkToken(kind) != true {
		fmt.Println(term_col_err, "Expected ", kind, " , got ", parser.curToken.kind.idname, term_col_res)
		os.Exit(1)
	}
	parser.nextToken()
}

// Check Token
func (parser *Parser) nextToken() {
	parser.curToken = parser.peekToken
	parser.peekToken = parser.lexer.getToken()
}

// Start parsing until EOF
// program ::= {statement}
func (parser *Parser) program() {
	// Emit C code with headers and main function
	emitheader.WriteString("#include <stdio.h>\n")
	emitheader.WriteString("int main(void){\n")

	// Since some newlines are required in our grammer,need to skip the access
	for parser.checkToken("NEWLINE") {
		parser.nextToken()
	}
	// Parse all the statements in the program
	for parser.checkToken("EOF") != true {
		parser.statement()
	}

	// Wrap things up.
	emitcode.WriteString("return 0;\n")
	emitcode.WriteString("}\n")

	// Check that each label referenced in a GOTO is declared.
	for _, lab := range parser.labelsGotoed {
		if sliceContains(parser.labelsDeclared, lab) != true {
			fmt.Println(term_col_err, "Attempting to GOTO to undeclared label: ", lab, term_col_res)
			os.Exit(1)
		}
	}
}

// Parse Statements
func (parser *Parser) statement() {
	// Check the first Token to see what kind of statement this is.

	if parser.checkToken("PRINT") { // "PRINT" (expression | string)
		parser.nextToken()

		if parser.checkToken("STRING") {
			// Simple String
			emitcode.WriteString("printf(\"" + parser.curToken.text + "\\n\");\n")
			parser.nextToken()
		} else {
			// Expect an expression and print the result as a float.
			emitcode.WriteString("printf(\"%" + ".2f\\n\",(float)(")
			// Expect an expression
			parser.expression()
			emitcode.WriteString("));\n")
		}
	} else if parser.checkToken("IF") { // "IF" comparison "THEN" {statement} "ENDIF"
		parser.nextToken()
		emitcode.WriteString("if(")
		parser.comparison()

		parser.matchToken("THEN")
		parser.nl()
		emitcode.WriteString("){")

		// Zero or more statements in the body.
		for parser.checkToken("ENDIF") != true {
			parser.statement()
		}

		parser.matchToken("ENDIF")
		emitcode.WriteString("}\n")
	} else if parser.checkToken("WHILE") { // "WHILE" comparison REPEAT" {statement} "ENDWHILE"
		parser.nextToken()
		emitcode.WriteString("while(")
		parser.comparison()

		parser.matchToken("REPEAT")
		parser.nl()
		emitcode.WriteString("){\n")

		// Zero or more statements in the loop body.
		for parser.checkToken("ENDWHILE") != true {
			parser.statement()
		}

		parser.matchToken("ENDWHILE")
		emitcode.WriteString("}\n")
	} else if parser.checkToken("LABEL") { // "LABEL" ident
		parser.nextToken()

		// Make sure this label doesn't already exist.
		if sliceContains(parser.labelsDeclared, parser.curToken.text) == true {
			fmt.Println(term_col_err, "LABEL already exists: ", parser.curToken.text, term_col_res)
			os.Exit(1)
		} else {
			parser.labelsDeclared = append(parser.labelsDeclared, parser.curToken.text)
		}
		emitcode.WriteString(parser.curToken.text + ":\n")
		parser.matchToken("IDENT")
	} else if parser.checkToken("GOTO") { // "GOTO" ident
		parser.nextToken()
		parser.labelsGotoed = append(parser.labelsGotoed, parser.curToken.text)
		emitcode.WriteString("goto" + parser.curToken.text + ";\n")
		parser.matchToken("IDENT")
	} else if parser.checkToken("LET") { // "LET" ident "=" expression
		parser.nextToken()
		// Check if ident exists in symbol table
		if sliceContains(parser.symbols, parser.curToken.text) != true {
			parser.symbols = append(parser.symbols, parser.curToken.text)
			emitheader.WriteString("float " + parser.curToken.text + ";\n")
		}
		emitcode.WriteString(parser.curToken.text + "=")
		parser.matchToken("IDENT")
		parser.matchToken("EQ")
		parser.expression()
		emitcode.WriteString(";\n")
	} else if parser.checkToken("INPUT") { // "INPUT" ident
		parser.nextToken()
		// Check if ident exists in symbol table
		if sliceContains(parser.symbols, parser.curToken.text) != true {
			parser.symbols = append(parser.symbols, parser.curToken.text)
			emitheader.WriteString("float " + parser.curToken.text + ";\n")
		}

		// Emit scanf but also validate the input. If invalid,set the variable to 0 and clear the input.
		emitcode.WriteString("if(0 == scanf(\"%" + "f\", &" + parser.curToken.text + ")) {\n")
		emitcode.WriteString(parser.curToken.text + " = 0;\n")
		emitcode.WriteString("scanf(\"%")
		emitcode.WriteString("*s\");\n")
		emitcode.WriteString("}\n")
		parser.matchToken("IDENT")
	} else { // This is not a valid statement. Fatal Error
		fmt.Println(term_col_err, "Invalid statement at ", parser.curToken.text, " ( ", parser.curToken.kind.idname, " ) ", term_col_res)
		os.Exit(1)
	}

	// Newline
	parser.nl()
}

// comparison ::= expression (("==" | "!=" | ">" | ">=" | "<" | "<=") expression)+
func (parser *Parser) comparison() {
	parser.expression()

	// Must be at least one comparison operator and another operation
	if parser.isComparisonOperator() == true {
		emitcode.WriteString(parser.curToken.text)
		parser.nextToken()
		parser.expression()
	} else {
		fmt.Println(term_col_err, "Expected comparison operator at:", parser.curToken.text, term_col_res)
		os.Exit(1)
	}

	// Can have 0 or more comparison operator and expressions.
	for parser.isComparisonOperator() == true {
		emitcode.WriteString(parser.curToken.text)
		parser.nextToken()
		parser.expression()
	}
}

// Check whether the token is comparison operator
func (parser *Parser) isComparisonOperator() bool {
	return (parser.checkToken("GT") || parser.checkToken("GTEQ") || parser.checkToken("LT") || parser.checkToken("LTEQ") || parser.checkToken("EQEQ") || parser.checkToken("NOTEQ"))
}

// expression ::= term {( "-" | "+" ) term}
func (parser *Parser) expression() {
	parser.term()

	// Can have 0 or more +/-  and expressions.
	for parser.checkToken("PLUS") || parser.checkToken("MINUS") {
		emitcode.WriteString(parser.curToken.text)
		parser.nextToken()
		parser.term()
	}
}

// term ::= unary {( "/" | "*" ) unary}
func (parser *Parser) term() {
	parser.unary()
	// Can have 0 or more *// and expressions.
	for parser.checkToken("ASTERISK") || parser.checkToken("SLASH") {
		emitcode.WriteString(parser.curToken.text)
		parser.nextToken()
		parser.unary()
	}
}

// unary ::= ["+" | "-"] primary
func (parser *Parser) unary() {
	// Optional Unary +/-
	if parser.checkToken("PLUS") || parser.checkToken("MINUS") {
		emitcode.WriteString(parser.curToken.text)
		parser.nextToken()
	}
	parser.primary()
}

// primary ::= number | ident
func (parser *Parser) primary() {
	if parser.checkToken("NUMBER") {
		emitcode.WriteString(parser.curToken.text)
		parser.nextToken()
	} else if parser.checkToken("IDENT") {
		// Ensure the variable already exists.
		if sliceContains(parser.symbols, parser.curToken.text) != true {
			fmt.Println(term_col_err, "Referencing variable before assignment: ", parser.curToken.text, term_col_res)
			os.Exit(1)
		}
		emitcode.WriteString(parser.curToken.text)
		parser.nextToken()
	} else { // unexpected token
		fmt.Println(term_col_err, "Unexpected Token at ", parser.curToken.text, term_col_res)
		os.Exit(1)
	}
}

// Newline
// nl ::= '\n'+
func (parser *Parser) nl() {
	// Require atleast one newline.
	parser.matchToken("NEWLINE")
	for parser.checkToken("NEWLINE") == true {
		parser.nextToken()
	}
}
