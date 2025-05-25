package main

import (
	"testing"
)

// BenchmarkLexer benchmarks the lexer performance
func BenchmarkLexer(b *testing.B) {
	input := `let x = 5;
	let y = 10;
	let add = func(a, b) {
		return a + b;
	};
	let result = add(x, y);
	if (result > 10) {
		return true;
	} else {
		return false;
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := New(input)
		for {
			tok := lexer.NextToken()
			if tok.Type == EOF {
				break
			}
		}
	}
}

// BenchmarkParser benchmarks the parser performance
func BenchmarkParser(b *testing.B) {
	input := `let x = 5;
	let y = 10;
	let add = func(a, b) {
		return a + b;
	};
	let result = add(x, y);
	if (result > 10) {
		return true;
	} else {
		return false;
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer := New(input)
		parser := NewParser(lexer)
		parser.ParseProgram()
	}
}

// BenchmarkEvaluator benchmarks the evaluator performance
func BenchmarkEvaluator(b *testing.B) {
	input := `let x = 5;
	let y = 10;
	let add = func(a, b) {
		return a + b;
	};
	let result = add(x, y);
	result;`

	lexer := New(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkFibonacci benchmarks a recursive function
func BenchmarkFibonacci(b *testing.B) {
	input := `func fib(n) {
		if (n <= 1) {
			return n;
		} else {
			return fib(n - 1) + fib(n - 2);
		}
	}
	fib(10);`

	lexer := New(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkComplexProgram benchmarks a complex program
func BenchmarkComplexProgram(b *testing.B) {
	input := `
	func factorial(n) {
		if (n <= 1) {
			return 1;
		} else {
			return n * factorial(n - 1);
		}
	}

	func sum(arr, index, total) {
		if (index >= 10) {
			return total;
		} else {
			return sum(arr, index + 1, total + index);
		}
	}

	let fact5 = factorial(5);
	let sum10 = sum([], 0, 0);
	let message = "Result: " + "factorial(5) = ";
	fact5;`

	lexer := New(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := NewEnvironment()
		Eval(program, env)
	}
}

// BenchmarkStringConcatenation benchmarks string operations
func BenchmarkStringConcatenation(b *testing.B) {
	input := `
	let a = "Hello";
	let b = " ";
	let c = "World";
	let d = "!";
	let result = a + b + c + d;
	result;`

	lexer := New(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := NewEnvironment()
		Eval(program, env)
	}
}
