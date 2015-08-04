package libloopcheck

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"reflect"
)

type debugging bool

// set to true to enable verbose debugging output
const debug debugging = false

func debugf(format string, args ...interface{}) {
	if debug {
		log.Printf(format, args...)
	}
}

func CheckFiles(paths []string) error {
	fset := token.NewFileSet()
	var pkgName string
	v := visitor{fset, bytes.Buffer{}, []warning{}}
	for _, filename := range paths {
		parsed, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			return err
		}
		if pkgName == "" {
			pkgName = parsed.Name.Name
		} else if parsed.Name.Name != pkgName {
			return fmt.Errorf("%s is in package %s, not %s", filename, parsed.Name.Name, pkgName)
		}

		ast.Walk(&v, parsed)
	}

	for _, warn := range v.warnings {
		usagePosition := fset.Position(warn.usagePos)

		fmt.Printf("%s:%d: takes address of loop variable: %s\n",
			usagePosition.Filename, usagePosition.Line, warn.usageExpr)
		rangePosition := fset.Position(warn.rangePos)
		fmt.Printf("  range at line %d: %s\n", rangePosition.Line, warn.rangeExpr)
	}
	return nil
}

type warning struct {
	usageExpr string
	usagePos  token.Pos
	rangeExpr string
	rangePos  token.Pos
}

type visitor struct {
	fset     *token.FileSet
	buffer   bytes.Buffer
	warnings []warning
}

// Inspired by go vet's rangeloop check:
// https://github.com/golang/tools/blob/master/cmd/vet/rangeloop.go
func (v *visitor) checkRange(node *ast.RangeStmt) bool {
	var keyObject *ast.Object
	var valueObject *ast.Object
	if key, ok := node.Key.(*ast.Ident); ok {
		keyObject = key.Obj
	} else if node.Key != nil {
		panic("wtf")
	}
	if value, ok := node.Value.(*ast.Ident); ok {
		valueObject = value.Obj
	} else if node.Value != nil {
		panic("wtf")
	}
	// no variables: no possible errors
	if keyObject == nil && valueObject == nil {
		return true
	}

	ast.Inspect(node.Body, func(n ast.Node) bool {
		if unary, ok := n.(*ast.UnaryExpr); ok {
			if unary.Op == token.AND {
				// recurse: could be address of a compound expression, e.g. &(rangeVar.y)
				// fmt.Printf("!! WTF %s %s\n", reflect.ValueOf(unary.X), v.str(unary.X))

				ast.Inspect(unary.X, func(n2 ast.Node) bool {
					// fmt.Printf("!!!! WTF %s %s\n", reflect.ValueOf(n2), v.str(n2))
					if ident, ok := n2.(*ast.Ident); ok {
						if ident.Obj == keyObject || ident.Obj == valueObject {

							v.warnings = append(v.warnings,
								warning{v.str(unary), unary.Pos(), v.rangestr(node), node.Pos()})
							// break early: no need to check more things
							return false
						}
					}
					return true
				})
				// we already recursively inspected the children
				return false
			}
		}
		return true
	})

	return true
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		// nil indicates the end of recursion; ignore it
		return v
	}

	debugf("%s\n", reflect.TypeOf(node))
	if ident, ok := node.(*ast.Ident); ok {
		debugf("  ident: %s\n", ident.String())
	} else if rs, ok := node.(*ast.RangeStmt); ok {
		debugf("  range: %s\n", v.rangestr(rs))
		ok = v.checkRange(rs)
		if !ok {
			fmt.Printf("#### RANGE LOOP ERROR\n")
		}
		// we've already checked the body of the range statement: don't recurse
		return nil
	}
	return v
}

func (v *visitor) rangestr(rs *ast.RangeStmt) string {
	v.buffer.Reset()
	v.buffer.WriteString("for ")
	printer.Fprint(&v.buffer, v.fset, rs.Key)
	if rs.Value != nil {
		v.buffer.WriteString(", ")
		printer.Fprint(&v.buffer, v.fset, rs.Value)
	}
	v.buffer.WriteString(" := range ")
	printer.Fprint(&v.buffer, v.fset, rs.X)
	return v.buffer.String()
}

func (v *visitor) str(node ast.Node) string {
	v.buffer.Reset()
	printer.Fprint(&v.buffer, v.fset, node)
	return v.buffer.String()
}
