package syntx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type RulePosition struct {
	ruleHash uint64
	position int
}

type ParseError struct {
	RuleName    string
	Position    int
	Description string
}

type Context struct {
	Rules         []int
	Literals      []string
	RuleJumpTable map[uint64]int
	RuleNameTable map[int]string
	JumpStack     *Stack[int]
	PositionStack *Stack[int]
	RootNode      *AstNode
	CurrentNode   *Stack[*AstNode]
	CurrentRule   *Stack[string]
	TextPosStack  *Stack[int]
	Error         ParseError
	FillTable     []RulePosition
}

func NewContext() Context {
	ctx := Context{
		RuleJumpTable: make(map[uint64]int),
		RuleNameTable: make(map[int]string),
		JumpStack:     NewStack[int](),
		PositionStack: NewStack[int](),
		RootNode:      NewAstNode("<root>", Range{0, 0}),
		CurrentNode:   NewStack[*AstNode](),
		CurrentRule:   NewStack[string](),
		TextPosStack:  NewStack[int](),
		Error:         ParseError{"", -1, ""},
	}

	ctx.CurrentNode.Push(ctx.RootNode)

	return ctx
}

func (ctx Context) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "Literals: %v\nRule names: %v\nRules:\n\t[\n", ctx.Literals, ctx.RuleNameTable)

	var nextCommandIndex int = 0

	for i, v := range ctx.Rules {
		fmt.Fprintf(&b, "\t\t%3d: %3d", i, v)

		if i == nextCommandIndex && 0 <= v && v < len(CommandNames) {
			fmt.Fprintf(&b, " (%s)", CommandNames[RulerType(v)])
			nextCommandIndex += ArgNums[RulerType(v)] + 1
		}

		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "\t]\n\n")

	fmt.Fprintf(&b, "AST:\n")

	printAst(&b, 1, ctx.RootNode)

	fmt.Fprintf(&b, "\n")
	printError(&b, ctx.Error)
	fmt.Fprintf(&b, "\n")

	return b.String()
}

func printAst(b *strings.Builder, tabs int, node *AstNode) {
	printTabs(b, tabs)
	fmt.Fprintf(b, "Name: %s (%p)\n", node.Name, node)

	printTabs(b, tabs)
	fmt.Fprintf(b, "Range: [%d, %d]\n", node.CoveredRange.From, node.CoveredRange.To)

	for _, child := range node.Children {
		printAst(b, tabs+1, child)
	}
}

func printTabs(b *strings.Builder, tabs int) {
	for i := 0; i < tabs; i++ {
		fmt.Fprintf(b, "\t")
	}
}

func printError(b *strings.Builder, parseError ParseError) {
	fmt.Fprintf(
		b, GetErrorMessage(parseError),
	)
}

func GetErrorMessage(parseError ParseError) string {
	if parseError.Position >= 0 {
		return fmt.Sprintf(
			"A %s rule failed while %s at position %d",
			parseError.RuleName, parseError.Description, parseError.Position,
		)
	} else {
		return ""
	}
}

func Unmarshal[T any](ast *AstNode, text string) (any, error) {
	dataType := reflect.TypeOf((*T)(nil)).Elem()

	instance := reflect.New(dataType)

	return unmarshalHelper(ast, instance, text)
}

func unmarshalHelper(ast *AstNode, data reflect.Value, text string) (any, error) {
	dataValue := reflect.Indirect(data)

	for _, child := range ast.Children {
		fieldName := child.Name
		fieldText := string(text[child.CoveredRange.From:child.CoveredRange.To])

		fmt.Printf("Looking for field name: %s\n", fieldName)
		var field reflect.Value = dataValue.FieldByName(fieldName)

		fmt.Printf("[Field] Type: \"%s\" (Kind: %s) Value: \"%v\" Valid: %t\n", field.Type(), field.Kind(), field, field.IsValid())

		switch field.Kind() {
		case reflect.Int:
			if value, err := strconv.ParseInt(fieldText, 10, 32); err != nil {
				return nil, err
			} else {
				fmt.Printf("Setting field to %d\n", value)
				field.SetInt(value)
			}
		case reflect.String:
			fmt.Printf("Setting field to %s\n", fieldText)
			field.SetString(fieldText)
		case reflect.Pointer:
			fmt.Printf("Pointer field: %s\n", field.Type())
			newInstance := reflect.New(field.Type().Elem())

			fmt.Printf(
				"New instance: %v %v %v %v\n",
				newInstance, newInstance.Elem(), newInstance.Pointer(), newInstance.UnsafePointer(),
			)
			field.Set(newInstance)

			if _, err := unmarshalHelper(child, newInstance, text); err != nil {
				return nil, fmt.Errorf("Error getting child's fields")
			}
		case reflect.Slice:
			fmt.Printf("Slice field: %s\n", field.Type())
			fmt.Printf("Slice's element type: %s\n", field.Type().Elem())
			elementType := field.Type().Elem()

			// Here we're expecting the slice to contain pointers to structs
			newInstance := reflect.New(elementType.Elem())
			field.Set(reflect.Append(field, newInstance))

			if _, err := unmarshalHelper(child, newInstance, text); err != nil {
				return nil, fmt.Errorf("Error getting child's fields")
			}
		}
	}

	return dataValue.Interface(), nil
}
