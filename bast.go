// Compiler for BAST(Basic Tiny) - A Tiny Toy Basic Dialect
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Use strings.Builder to properly concatenate strings(required by emitter)
var (
	emitheader strings.Builder
	emitcode   strings.Builder
)

// color strings required to color text in terminal
const (
	term_col_err string = "\u001b[31m" // Red color
	term_col_suc string = "\u001b[32m" // Green color
	term_col_res string = "\u001b[0m"  // color reset
)

// Main Function.Parse flags,Initialize lexer+parser and compile
func main() {
	fmt.Println("Compiler for BAST(BASIC Tiny) - A Tiny Toy Basic Dialect (--help/-h for more help)")
	// Specify Command Line Flags
	pathin := flag.String("in", "inputfile", "Path to Input Source File")
	sourceout := flag.String("out", "outputfile.c", "Path to Output Emitted C")
	flag.Parse()

	input, err := ioutil.ReadFile(*pathin)
	if err != nil {
		fmt.Fprintf(os.Stderr, term_col_err+"fatal error: %v\ncompilation terminated.\nspecify input file by -in=inputfile.bast\n"+term_col_res, err)
		os.Exit(1)
	}
	var sourcein string = string(input)

	// Initialize the lexer , emitter and parser.
	// Lexer
	var lexerinst Lexer
	lexerinst.initlexer(sourcein)
	// Parser
	var parserinst Parser
	parserinst.initParser(lexerinst)

	// Start Parsing
	parserinst.program()
	// Write the output to file
	writebuf := []byte(emitheader.String() + emitcode.String())
	fmt.Print("Writing to the File >> ", "\u001b[33m", *sourceout, term_col_res)
	err = ioutil.WriteFile(*sourceout, writebuf, 0644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v", err)
		os.Exit(1)
	}

	// Notify Completion of compilation
	fmt.Print(term_col_suc, "\nCompiling Completed...\nYou can now complile emitted C code ", "\u001b[33m", *sourceout, term_col_suc, " with your favorite C compiler.\n", term_col_res)
}
