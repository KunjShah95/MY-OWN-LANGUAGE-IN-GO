package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Token structure
type Token struct {
	Type  string
	Value string
}

// Lexer function
func lexer(input string) []Token {
	var tokens []Token
	cursor := 0

	for cursor < len(input) {
		char := string(input[cursor])

		// Skip whitespace
		if match, _ := regexp.MatchString(`\s`, char); match {
			cursor++
			continue
		}

		// Check for characters
		if match, _ := regexp.MatchString(`[a-zA-Z]`, char); match {
			word := ""
			for cursor < len(input) && regexp.MustCompile(`[a-zA-Z]`).MatchString(string(input[cursor])) {
				word += string(input[cursor])
				cursor++
			}
			if word == "ye" || word == "bol" || word == "agar" || word == "warna" || word == "jabtak" || word == "switch" || word == "case" || word == "default" || word == "array" || word == "file" {
				tokens = append(tokens, Token{Type: "keyword", Value: word})
			} else {
				tokens = append(tokens, Token{Type: "identifier", Value: word})
			}
			continue
		}

		// Check for numbers
		if match, _ := regexp.MatchString(`[0-9]`, char); match {
			number := ""
			for cursor < len(input) && regexp.MustCompile(`[0-9]`).MatchString(string(input[cursor])) {
				number += string(input[cursor])
				cursor++
			}
			tokens = append(tokens, Token{Type: "number", Value: number})
			continue
		}

		// Tokenize operators and equals sign
		if match, _ := regexp.MatchString(`[\+\-\*\/=<>]`, char); match {
			tokens = append(tokens, Token{Type: "operator", Value: char})
			cursor++
			continue
		}
	}

	return tokens
}

// AST Node structures
type Node interface{}

type Program struct {
	Body []Node
}

type Declaration struct {
	Name  string
	Value int
}

type Print struct {
	Expression string
}

type Conditional struct {
	Condition string
	Body      []Node
	ElseBody  []Node
}

type Loop struct {
	Condition string
	Body      []Node
}

type SwitchCase struct {
	Expression string
	Cases      map[string][]Node
	Default    []Node
}

type Array struct {
	Name   string
	Values []int
}

type FileOperation struct {
	Operation string
	Filename  string
	Content   string
}

// Parser function
func parser(tokens []Token) Program {
	var variables = map[string]int{}
	program := Program{Body: []Node{}}

	for len(tokens) > 0 {
		token := tokens[0]
		tokens = tokens[1:]

		// Handle declarations
		if token.Type == "keyword" && token.Value == "ye" {
			declaration := Declaration{Name: tokens[0].Value}
			tokens = tokens[1:]

			// Check for assignment
			if len(tokens) > 0 && tokens[0].Type == "operator" && tokens[0].Value == "=" {
				tokens = tokens[1:] // Consume '='
				expression := ""
				for len(tokens) > 0 && tokens[0].Type != "keyword" {
					expression += tokens[0].Value
					tokens = tokens[1:]
				}

				// Evaluate the expression
				declaration.Value = evaluateExpression(expression, variables)
				variables[declaration.Name] = declaration.Value
			}

			program.Body = append(program.Body, declaration)
		}

		// Handle print statements
		if token.Type == "keyword" && token.Value == "bol" {
			printNode := Print{Expression: tokens[0].Value}
			tokens = tokens[1:]
			program.Body = append(program.Body, printNode)
		}

		// Handle conditionals
		if token.Type == "keyword" && token.Value == "agar" {
			condition := tokens[0].Value
			tokens = tokens[1:]
			body := []Node{}
			elseBody := []Node{}

			// Parse the body of the conditional
			for len(tokens) > 0 && tokens[0].Value != "warna" && tokens[0].Value != "end" {
				body = append(body, parser(tokens).Body...)
			}

			// Handle else part
			if len(tokens) > 0 && tokens[0].Value == "warna" {
				tokens = tokens[1:]
				for len(tokens) > 0 && tokens[0].Value != "end" {
					elseBody = append(elseBody, parser(tokens).Body...)
				}
			}

			program.Body = append(program.Body, Conditional{Condition: condition, Body: body, ElseBody: elseBody})
		}

		// Handle loops
		if token.Type == "keyword" && token.Value == "jabtak" {
			condition := tokens[0].Value
			tokens = tokens[1:]
			body := []Node{}

			// Parse the body of the loop
			for len(tokens) > 0 && tokens[0].Value != "end" {
				body = append(body, parser(tokens).Body...)
			}

			program.Body = append(program.Body, Loop{Condition: condition, Body: body})
		}

		// Handle switch cases
		if token.Type == "keyword" && token.Value == "switch" {
			expression := tokens[0].Value
			tokens = tokens[1:]
			cases := map[string][]Node{}
			defaultBody := []Node{}

			// Parse the cases
			for len(tokens) > 0 && tokens[0].Value != "end" {
				if tokens[0].Value == "case" {
					caseValue := tokens[1].Value
					tokens = tokens[2:]
					caseBody := []Node{}
					for len(tokens) > 0 && tokens[0].Value != "case" && tokens[0].Value != "default" && tokens[0].Value != "end" {
						caseBody = append(caseBody, parser(tokens).Body...)
					}
					cases[caseValue] = caseBody
				} else if tokens[0].Value == "default" {
					tokens = tokens[1:]
					for len(tokens) > 0 && tokens[0].Value != "end" {
						defaultBody = append(defaultBody, parser(tokens).Body...)
					}
				}
			}

			program.Body = append(program.Body, SwitchCase{Expression: expression, Cases: cases, Default: defaultBody})
		}

		// Handle arrays
		if token.Type == "keyword" && token.Value == "array" {
			name := tokens[0].Value
			tokens = tokens[1:]
			values := []int{}

			// Parse array values
			for len(tokens) > 0 && tokens[0].Type == "number" {
				value, _ := strconv.Atoi(tokens[0].Value)
				values = append(values, value)
				tokens = tokens[1:]
			}

			program.Body = append(program.Body, Array{Name: name, Values: values})
		}

		// Handle file operations
		if token.Type == "keyword" && token.Value == "file" {
			operation := tokens[0].Value
			filename := tokens[1].Value
			content := ""
			if operation == "write" {
				content = tokens[2].Value
				tokens = tokens[3:]
			} else {
				tokens = tokens[2:]
			}

			program.Body = append(program.Body, FileOperation{Operation: operation, Filename: filename, Content: content})
		}
	}

	return program
}

// Evaluate expressions
func evaluateExpression(expr string, vars map[string]int) int {
	// Replace variable names with their values
	for key, value := range vars {
		expr = strings.ReplaceAll(expr, key, strconv.Itoa(value))
	}

	// Evaluate the expression using strconv and a simple arithmetic parser
	result, _ := eval(expr)
	return result
}

// Simple arithmetic expression evaluator
func eval(expr string) (int, error) {
	// Remove all spaces
	expr = strings.ReplaceAll(expr, " ", "")

	// Parse and evaluate the expression
	return parseExpr(expr)
}

func parseExpr(expr string) (int, error) {
	// Handle addition and subtraction
	for i := len(expr) - 1; i >= 0; i-- {
		if expr[i] == '+' || expr[i] == '-' {
			left, err := parseExpr(expr[:i])
			if err != nil {
				return 0, err
			}
			right, err := parseTerm(expr[i+1:])
			if err != nil {
				return 0, err
			}
			if expr[i] == '+' {
				return left + right, nil
			}
			return left - right, nil
		}
	}
	return parseTerm(expr)
}

func parseTerm(expr string) (int, error) {
	// Handle multiplication and division
	for i := len(expr) - 1; i >= 0; i-- {
		if expr[i] == '*' || expr[i] == '/' {
			left, err := parseTerm(expr[:i])
			if err != nil {
				return 0, err
			}
			right, err := parseFactor(expr[i+1:])
			if err != nil {
				return 0, err
			}
			if expr[i] == '*' {
				return left * right, nil
			}
			return left / right, nil
		}
	}
	return parseFactor(expr)
}

func parseFactor(expr string) (int, error) {
	// Handle numbers
	return strconv.Atoi(expr)
}

// Code generator function
func codeGen(node Node) string {
	switch n := node.(type) {
	case Program:
		var code []string
		for _, child := range n.Body {
			code = append(code, codeGen(child))
		}
		return strings.Join(code, "\n")
	case Declaration:
		return fmt.Sprintf("const %s = %d;", n.Name, n.Value)
	case Print:
		return fmt.Sprintf("fmt.Println(%s)", n.Expression)
	case Conditional:
		condition := fmt.Sprintf("if %s {", n.Condition)
		body := codeGen(Program{Body: n.Body})
		elseBody := ""
		if len(n.ElseBody) > 0 {
			elseBody = fmt.Sprintf("} else {%s", codeGen(Program{Body: n.ElseBody}))
		}
		return fmt.Sprintf("%s\n%s\n%s}", condition, body, elseBody)
	case Loop:
		condition := fmt.Sprintf("for %s {", n.Condition)
		body := codeGen(Program{Body: n.Body})
		return fmt.Sprintf("%s\n%s\n}", condition, body)
	case SwitchCase:
		expression := fmt.Sprintf("switch %s {", n.Expression)
		cases := ""
		for caseValue, caseBody := range n.Cases {
			cases += fmt.Sprintf("case %s:\n%s\n", caseValue, codeGen(Program{Body: caseBody}))
		}
		defaultBody := ""
		if len(n.Default) > 0 {
			defaultBody = fmt.Sprintf("default:\n%s\n", codeGen(Program{Body: n.Default}))
		}
		return fmt.Sprintf("%s\n%s\n%s}", expression, cases, defaultBody)
	case Array:
		values := ""
		for _, value := range n.Values {
			values += fmt.Sprintf("%d, ", value)
		}
		return fmt.Sprintf("var %s = []int{%s};", n.Name, values)
	case FileOperation:
		if n.Operation == "write" {
			return fmt.Sprintf("os.WriteFile(%s, []byte(%s), 0644)", n.Filename, n.Content)
		} else if n.Operation == "read" {
			return fmt.Sprintf("content, _ := os.ReadFile(%s)\nfmt.Println(string(content))", n.Filename)
		}
	}
	return ""
}

// Compiler function
func compiler(input string) string {
	tokens := lexer(input)
	ast := parser(tokens)
	executableCode := codeGen(ast)
	return executableCode
}

func main() {
	code := `
ye x = 10
ye y = 20
agar x < y
    bol x
warna
    bol y
end

jabtak x < y
    bol x
    ye x = x + 1
end

switch x
case 10
    bol "x is 10"
case 20
    bol "x is 20"
default
    bol "x is something else"
end

array arr 1 2 3 4 5
bol arr

file write "test.txt" "Hello, World!"
file read "test.txt"
`

	// Generate Go code
	tokens := lexer(code)
	ast := parser(tokens)

	// Define variables map
	variables := map[string]int{}

	// Evaluate the result for the `bol` statement
	for _, node := range ast.Body {
		switch n := node.(type) {
		case Declaration:
			variables[n.Name] = n.Value
		case Print:
			fmt.Println("Output:", evaluateExpression(n.Expression, variables)) // Print evaluated result
		case Conditional:
			// Handle conditional statements
		case Loop:
			// Handle loops
		case SwitchCase:
			// Handle switch cases
		case Array:
			// Handle arrays
		case FileOperation:
			// Handle file operations
		}
	}
}
