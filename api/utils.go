package api

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/iancoleman/orderedmap"
)

// parseGoExpression takes a Go expression in string form and evaluates it into an actual Go value.
func ParseGoExpression(expr string) (interface{}, error) {
	// Use the Go parser to validate and parse the expression
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression: %w", err)
	}

	// Handle specific cases based on the expression type
	switch v := node.(type) {
	case *ast.CompositeLit: // Handle arrays/slices
		return ParseCompositeLit(v)
	case *ast.BasicLit: // Handle basic literals
		return ParseBasicLit(v)
	case *ast.UnaryExpr: // Handle unary expressions (e.g., -5, +0)
		return ParseUnaryExpr(v)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", v)
	}
}

// parseCompositeLit parses a Go composite literal (e.g., `[]int{1, 2, 3}`)
func ParseCompositeLit(lit *ast.CompositeLit) (interface{}, error) {
	// Example assumes array literals of type []int
	var result []int
	for _, elt := range lit.Elts {
		if basicLit, ok := elt.(*ast.BasicLit); ok {
			var value int
			fmt.Sscanf(basicLit.Value, "%d", &value) // Basic conversion from string
			result = append(result, value)
		} else {
			return nil, fmt.Errorf("unsupported composite element type: %T", elt)
		}
	}
	return result, nil
}

// parseBasicLit parses a Go basic literal (e.g., "5")
func ParseBasicLit(lit *ast.BasicLit) (interface{}, error) {
	switch lit.Kind {
	case token.INT: // Integer
		var value int
		fmt.Sscanf(lit.Value, "%d", &value)
		return value, nil
	case token.STRING: // String
		return lit.Value[1 : len(lit.Value)-1], nil // Remove quotes
	default:
		return nil, fmt.Errorf("unsupported literal kind: %s", lit.Kind)
	}
}

// parseUnaryExpr parses a Go unary expression (e.g., -5, +0)
func ParseUnaryExpr(expr *ast.UnaryExpr) (interface{}, error) {
	// Only handle basic literals as the operand
	if basicLit, ok := expr.X.(*ast.BasicLit); ok {
		value, err := ParseBasicLit(basicLit)
		if err != nil {
			return nil, err
		}

		// Apply the unary operator
		switch expr.Op {
		case token.SUB: // Negative numbers
			if intValue, ok := value.(int); ok {
				return -intValue, nil
			}
		case token.ADD: // Positive numbers (no-op)
			return value, nil
		default:
			return nil, fmt.Errorf("unsupported unary operator: %s", expr.Op)
		}
	}

	return nil, fmt.Errorf("unsupported operand type for unary expression: %T", expr.X)
}

// Unmarshal JSON and format for unit test execution
func FormatTestJSON(inputJSON string) string {
	o := orderedmap.New()
	err := json.Unmarshal([]byte(inputJSON), &o)
	if err != nil {
		log.Fatalf("Failed to parse input: %v", err)
	}

	keys := o.Keys()
	outputString, sep := "", ""
	for _, k := range keys {
		v, _ := o.Get(k)
		s, ok := v.(string)
		if ok {
			outputString += sep + s
			sep = ", "
		} else {
			log.Fatalf("JSON value is not a string.")
		}
	}

	outputString = strings.TrimSuffix(outputString, ", ")
	outputString = NormalizeString(outputString)

	return outputString
}

// Regex to remove leading or trailing backslashes
func RemoveOuterBackslashes(s string) string {
	re := regexp.MustCompile(`^\\|\\$`)
	return re.ReplaceAllString(s, "")
}

// Replace escaped quotes (\") with actual quotes (")
func NormalizeString(s string) string {
	s = strings.ReplaceAll(s, `\"`, `"`)
	return s
}

func MapType(typeName string) reflect.Type {
	switch typeName {
	case "int":
		return reflect.TypeOf(0)
	case "string":
		return reflect.TypeOf("")
	case "bool":
		return reflect.TypeOf(false)
	case "[]int":
		return reflect.TypeOf([]int{})
	// Add other types as needed
	default:
		return nil // Unknown type
	}
}

func BuildFuncType(inputTypes []string, outputTypes []string) reflect.Type {
	var in []reflect.Type
	for _, t := range inputTypes {
		mappedType := MapType(t)
		if mappedType == nil {
			log.Fatalf("Unsupported input type: %s", t)
		}
		in = append(in, mappedType)
	}

	var out []reflect.Type
	for _, t := range outputTypes {
		mappedType := MapType(t)
		if mappedType == nil {
			log.Fatalf("Unsupported output type: %s", t)
		}
		out = append(out, mappedType)
	}

	return reflect.FuncOf(in, out, false)
}
