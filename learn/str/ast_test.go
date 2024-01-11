package str

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Go 代码
var code = `
package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

const AA = "aa"

type Customer struct {
		Id   int32  
		Name string 
}
func (c *Customer) hello() {

}

func main() {
	fmt.Println("Hello, AST!")
}
`

// TestAst 测试ast的方法
//
//	ast用于读取go文件，根据文件的定义，可以获取go定义信息，可以结合template生成新的go文件
func TestAst(t *testing.T) {
	node, err := parser.ParseFile(token.NewFileSet(), "demo.go", code, parser.ParseComments)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf(" - pageage:  %s\n", node.Name.Name) //函数

	for _, d := range node.Decls {
		switch decl := d.(type) {
		case *ast.FuncDecl:
			fmt.Printf("Function Declaration: %s\n", decl.Name.Name) //函数
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch s := spec.(type) {
				case *ast.ImportSpec:
					fmt.Printf(" - Import: %s\n", s.Path.Value) //import
				case *ast.ValueSpec:
					fmt.Printf(" - Value: %s\n", s.Names[0].Name) //const
				case *ast.TypeSpec:
					fmt.Printf(" - Type: %s\n", s.Name.Name) //struct
				}
			}
		}
	}
}
