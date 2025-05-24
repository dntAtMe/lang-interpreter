package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tinylang <file.tiny>")
		return
	}

	filename := os.Args[1]
	if !strings.HasSuffix(filename, ".tiny") {
		fmt.Println("Error: File must have .tiny extension")
		return
	}

	RunFile(filename)
}

// getStatementType returns a human-readable description of the statement type
func getStatementType(stmt Statement) string {
	switch s := stmt.(type) {
	case *LetStatement:
		return fmt.Sprintf("Let Statement (var: %s)", s.Name.Value)
	case *FunctionStatement:
		return fmt.Sprintf("Function Statement (name: %s, params: %d)", s.Name.Value, len(s.Parameters))
	case *ReturnStatement:
		if s.ReturnValue != nil {
			return "Return Statement (with value)"
		}
		return "Return Statement (void)"
	case *IfStatement:
		if s.Alternative != nil {
			return "If-Else Statement"
		}
		return "If Statement"
	case *ExpressionStatement:
		return fmt.Sprintf("Expression Statement (%s)", getExpressionType(s.Expression))
	case *BlockStatement:
		return fmt.Sprintf("Block Statement (%d statements)", len(s.Statements))
	default:
		return "Unknown Statement"
	}
}

// getExpressionType returns a human-readable description of the expression type
func getExpressionType(expr Expression) string {
	switch e := expr.(type) {
	case *Identifier:
		return fmt.Sprintf("Identifier (%s)", e.Value)
	case *IntegerLiteral:
		return fmt.Sprintf("Integer (%d)", e.Value)
	case *StringLiteral:
		return fmt.Sprintf("String (%q)", e.Value)
	case *BooleanLiteral:
		return fmt.Sprintf("Boolean (%t)", e.Value)
	case *PrefixExpression:
		return fmt.Sprintf("Prefix (%s)", e.Operator)
	case *InfixExpression:
		return fmt.Sprintf("Infix (%s)", e.Operator)
	case *CallExpression:
		return fmt.Sprintf("Call (%d args)", len(e.Arguments))
	default:
		return "Unknown Expression"
	}
}
