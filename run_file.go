package main

import (
	"fmt"
	"os"
)

// RunFile executes a TinyLang file
func RunFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	code := string(content)
	lexer := New(code)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		fmt.Println("Parse errors:")
		for _, err := range parser.Errors() {
			fmt.Printf("  - %s\n", err)
		}
		return
	}

	env := NewEnvironment()
	result := Eval(program, env)

	if result.Type() == ERROR_OBJ {
		fmt.Printf("Runtime Error: %s\n", result.Inspect())
		return
	}

	fmt.Println(result.Inspect())
}
