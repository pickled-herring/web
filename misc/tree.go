package main

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

func min(a,b uint32)uint32{
	if a < b {
		return a
	}
	return b
}
func main() {
	input := []byte("function hello() { console.log('hello') }; function goodbye(){}")

	fmt.Println("len input:", len(input))
	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	tree := parser.Parse(nil, input)

	n := tree.RootNode()

	fmt.Println("AST:", n)
	fmt.Println("Root type:", n.Type())
	fmt.Println("Root children:", n.ChildCount())

	fmt.Println("\nFunctions in input:")
	q, _ := sitter.NewQuery([]byte("(function_declaration) @func"), javascript.GetLanguage())
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	var funcs []*sitter.Node
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			funcs = append(funcs, c.Node)
			fmt.Println("-", funcName(input, c.Node))
		}
	}

	fmt.Println("\nEdit input")
	input = []byte("function hello() { console.log('hello') }; function goodbye(){ console.log('goodbye') }")
	read := func(offset uint32, position sitter.Point) []byte {
		fmt.Printf("called read func with offset = %d and position = %#v\n", offset, position)
		if position.Row != 0 {
			return nil
		}
		return input[position.Column:min(position.Column+3, uint32(len(input)))]
	}
	input_change := sitter.Input{
		Read: read,
		Encoding: sitter.InputEncodingUTF8,
	}
	// reuse tree
	tree.Edit(sitter.EditInput{
		StartIndex:  62,
		OldEndIndex: 63,
		NewEndIndex: 87,
		StartPoint: sitter.Point{
			Row:    0,
			Column: 62,
		},
		OldEndPoint: sitter.Point{
			Row:    0,
			Column: 63,
		},
		NewEndPoint: sitter.Point{
			Row:    0,
			Column: 87,
		},
	})

	for _, f := range funcs {
		var textChange string
		if f.HasChanges() {
			textChange = "has change"
		} else {
			textChange = "no changes"
		}
		fmt.Println("-", funcName(input, f), ">", textChange)
	}

	newTree := parser.ParseInput(tree, input_change)
	n = newTree.RootNode()
	fmt.Println("\nNew AST:", n)
}

func funcName(content []byte, n *sitter.Node) string {
	if n == nil {
		return ""
	}

	if n.Type() != "function_declaration" {
		return ""
	}

	return n.ChildByFieldName("name").Content(content)
}
