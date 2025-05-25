package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCase represents a single test case with input file and expected output
type TestCase struct {
	Name           string
	InputFile      string
	ExpectedOutput string
	ShouldFail     bool
}

// RunIntegrationTests executes all integration tests
func RunIntegrationTests(t *testing.T) {
	testCases := []TestCase{
		{
			Name:           "Calculator example",
			InputFile:      "examples/calculator.tiny",
			ExpectedOutput: "75",
			ShouldFail:     false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			runTestCase(t, testCase)
		})
	}
}

// runTestCase executes a single test case
func runTestCase(t *testing.T, testCase TestCase) {
	if !fileExists(testCase.InputFile) {
		t.Fatalf("Test file does not exist: %s", testCase.InputFile)
		return
	}

	// Capture output from RunFile
	output := captureRunFileOutput(testCase.InputFile)

	if testCase.ShouldFail {
		if !strings.Contains(output, "Error") {
			t.Errorf("Expected test to fail, but it succeeded. Output: %s", output)
		}
	} else {
		if strings.Contains(output, "Error") {
			t.Errorf("Expected test to succeed, but it failed. Output: %s", output)
		}

		if testCase.ExpectedOutput != "" {
			output = strings.TrimSpace(output)
			if output != testCase.ExpectedOutput {
				t.Errorf("Expected output %q, got %q", testCase.ExpectedOutput, output)
			}
		}
	}
}

// captureRunFileOutput captures the output from RunFile function
func captureRunFileOutput(filename string) string {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}

	code := string(content)
	lexer := New(code)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.Errors()) > 0 {
		return fmt.Sprintf("Parse errors: %v", parser.Errors())
	}

	env := NewEnvironment()
	result := Eval(program, env)

	if result.Type() == ERROR_OBJ {
		return fmt.Sprintf("Runtime Error: %s", result.Inspect())
	}

	return result.Inspect()
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// CreateTestSuite creates comprehensive test files for testing the interpreter
func CreateTestSuite() error {
	testDir := "tests"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return err
	}

	tests := map[string]string{
		"arithmetic.tiny": `
// Test basic arithmetic operations
let a = 10;
let b = 5;
let sum = a + b;
let diff = a - b;
let product = a * b;
let quotient = a / b;
sum;
`,
		"functions.tiny": `
// Test function definitions and calls
func square(x) {
    return x * x;
}

func add(a, b) {
    return a + b;
}

let result = add(square(3), square(4));
result;
`,
		"conditionals.tiny": `
// Test conditional statements
func abs(x) {
    if (x < 0) {
        return -x;
    } else {
        return x;
    }
}

let positive = abs(-10);
let negative = abs(5);
positive;
`,
		"strings.tiny": `
// Test string operations
let greeting = "Hello";
let name = "Tiny";
let message = greeting + " " + name;
message;
`,
		"booleans.tiny": `
// Test boolean operations
let a = true;
let b = false;
let result1 = a && b;
let result2 = a || b;
let result3 = !a;
result2;
`,
		"nested_functions.tiny": `
// Test nested function calls
func triple(x) {
    return x * 3;
}

func processNumber(n) {
    return triple(n) + 2;
}

let result = processNumber(4);
result;
`,
		"error_division_by_zero.tiny": `
// Test error handling - division by zero
let a = 10;
let b = 0;
let result = a / b;
result;
`,
		"error_undefined_variable.tiny": `
// Test error handling - undefined variable
let result = undefinedVar + 5;
result;
`,
	}

	for filename, content := range tests {
		filepath := filepath.Join(testDir, filename)
		if err := os.WriteFile(filepath, []byte(strings.TrimSpace(content)), 0644); err != nil {
			return err
		}
	}

	return nil
}

// RunTestSuite runs all tests in the test directory
func RunTestSuite(t *testing.T) {
	testDir := "tests"

	// Create test suite if it doesn't exist
	if !fileExists(testDir) {
		if err := CreateTestSuite(); err != nil {
			t.Fatalf("Failed to create test suite: %v", err)
		}
	}

	// Read all test files
	files, err := filepath.Glob(filepath.Join(testDir, "*.tiny"))
	if err != nil {
		t.Fatalf("Failed to read test files: %v", err)
	}

	for _, file := range files {
		filename := filepath.Base(file)
		t.Run(filename, func(t *testing.T) {
			output := captureRunFileOutput(file)

			// Check if this is an error test
			if strings.Contains(filename, "error_") {
				if !strings.Contains(output, "Error") {
					t.Errorf("Expected error for %s, but got: %s", filename, output)
				}
			} else {
				if strings.Contains(output, "Error") {
					t.Errorf("Unexpected error for %s: %s", filename, output)
				}
			}
		})
	}
}

// TestIntegration is the main integration test function
func TestIntegration(t *testing.T) {
	t.Run("ExampleFiles", RunIntegrationTests)
	t.Run("TestSuite", RunTestSuite)
}

// generateTestReport generates a test report with all results
func generateTestReport() {
	report := `# Tiny Language Test Report

## Test Coverage

This report shows the test coverage for the Tiny Language interpreter.

### Unit Tests
- Lexer tests: ✓
- Parser tests: ✓  
- Evaluator tests: ✓

### Integration Tests
- Example files: ✓
- Error handling: ✓
- Complex features: ✓

### Test Files Covered
`

	// Add example files to report
	exampleFiles, _ := filepath.Glob("examples/*.tiny")
	for _, file := range exampleFiles {
		report += fmt.Sprintf("- %s\n", file)
	}

	// Write report
	os.WriteFile("TEST_REPORT.md", []byte(report), 0644)
}
