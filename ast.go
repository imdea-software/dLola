package dLola

import (
	//	"errors"
	"fmt"
	//	"log"
	//	"strconv"
	//	"strings"
)

type StreamType int

const (
	NumT StreamType = iota
	BoolT
	StringT
	Unknown  // we use this in the parser for unknow type values (offset expressions) that will be resolved later
	LastType = StringT
)

type StreamName string

func (s StreamName) Sprint() string {
	return string(s)
}

func (t StreamType) Sprint() string {

	type_names := []string{"num", "bool", "string"}

	// str string
	// switch t {
	// case Int:
	// 	str = "int"
	// case Bool:
	// 	str = "bool"
	// case String:
	// 	str = "string"
	// }
	// return str

	if t >= LastType {
		return ""
	}
	return fmt.Sprintf("%s", type_names[t])
}

type ConstDecl struct { // const int one_sec := 1s
	Name StreamName
	Type StreamType
	Val  Expr
}
type InputDecl struct { // input int bar
	Name StreamName
	Type StreamType
}
type OutputDecl struct { // output int foo /* this is just a decl, later a tick and a def will be given */
	Name StreamName
	Type StreamType
}

/*lm: not needed type TicksDecl struct {
	Name  StreamName
	Ticks TickingExpr
}*/

type OutputDefinition struct {
	Name StreamName
	Type StreamType
	Expr Expr // chango to ValueExpr?
}

func NewConstDecl(n, t, e interface{}) ConstDecl {
	name := getStreamName(n)
	return ConstDecl{name, t.(StreamType), e.(Expr)}
}
func NewInputDecl(n, t interface{}) InputDecl {
	name := getStreamName(n)
	return InputDecl{name, t.(StreamType)}
}
func NewOutputDecl(n, t interface{}) OutputDecl {
	name := getStreamName(n)
	return OutputDecl{name, t.(StreamType)}
}

/*: not needed func NewTicksDecl(n, t interface{}) TicksDecl {
	name := getStreamName(n)
	expr := t.(TickingExpr)
	return TicksDecl{name, expr}
}*/

func NewOutputDefinition(n, t, e interface{}) OutputDefinition {
	name := getStreamName(n)
	expr := e.(Expr)
	return OutputDefinition{name, t.(StreamType), expr}
}

func getStreamName(a interface{}) StreamName {
	return StreamName(a.(Identifier).Val)
}

func (c ConstDecl) Sprint() string {
	return fmt.Sprintf("ConstDecl: {Name = %s, Type = %s, expr = %s}", c.Name.Sprint(), c.Type.Sprint(), c.Val.Sprint())
}

func (i InputDecl) Sprint() string {
	return fmt.Sprintf("InputDecl: {Name = %s, Type = %s}", i.Name.Sprint(), i.Type.Sprint())
}

func (o OutputDecl) Sprint() string {
	return fmt.Sprintf("OutputDecl: {Name = %s, Type = %s}", o.Name.Sprint(), o.Type.Sprint())
}

func (o OutputDefinition) Sprint() string {
	return fmt.Sprintf("OutputDefinition: {Name = %s, Type = %s, expr = %s}", o.Name.Sprint(), o.Type.Sprint(), o.Expr.Sprint())
}

/*Pretty Print using method Accept*/
func (c ConstDecl) PrettyPrint() string {
	v := PrettyPrinterVisitor{0, ""}
	c.Val.Accept(&v)
	return fmt.Sprintf("ConstDecl: {Name = %s, Type = %s, expr = %s}", c.Name.Sprint(), c.Type.Sprint(), v.s)
}

func (i InputDecl) PrettyPrint() string {
	return fmt.Sprintf("InputDecl: {Name = %s, Type = %s}", i.Name.Sprint(), i.Type.Sprint())
}

func (o OutputDecl) PrettyPrint() string {
	return fmt.Sprintf("OutputDecl: {Name = %s, Type = %s}", o.Name.Sprint(), o.Type.Sprint())
}

func (o OutputDefinition) PrettyPrint() string {
	v := PrettyPrinterVisitor{0, ""}
	o.Expr.Accept(&v)
	return fmt.Sprintf("OutputDefinition: {Name = %s, Type = %s, expr = \n%s}", o.Name.Sprint(), o.Type.Sprint(), v.s)
}

/*Calls to the parser generated by PIGEON and returns a list of all the trees matched (each sentence will have a tree)*/
func GetAst(filename, prefix string) []interface{} {
	ast, err := ParseFile(filename)
	if err != nil {
		fmt.Printf(prefix+"There was an error: %s\n", err)
		return []interface{}{}
	}
	last := ast.([]interface{})
	/*TODO: perform castings before returning -> it wont work since the output type will still be []interface{}, so direct call to methods won't work*/
	return last
}

/*Gets the output of GetAst and prints it*/
func PrintAst(ast []interface{}, prefix string) {
	for _, val := range ast {
		switch v := val.(type) {
		/*	case dLola.Spec:
			fmt.Printf(prefix+"AST spec %s\n", v.Sprint())*/
		case ConstDecl:
			fmt.Printf(prefix+"%s\n", v.Sprint())
		case InputDecl:
			fmt.Printf(prefix+"%s\n", v.Sprint())
		case OutputDecl:
			fmt.Printf(prefix+"%s\n", v.Sprint())
		case OutputDefinition:
			fmt.Printf(prefix+"%s\n", v.Sprint())

		}
	}

}

func PrettyPrintAst(ast []interface{}, prefix string) {
	for _, val := range ast {
		switch v := val.(type) {
		/*	case dLola.Spec:
			fmt.Printf(prefix+"AST spec %s\n", v.Sprint())*/
		case ConstDecl:
			fmt.Printf(prefix+"%s\n", v.PrettyPrint())
		case InputDecl:
			fmt.Printf(prefix+"%s\n", v.PrettyPrint())
		case OutputDecl:
			fmt.Printf(prefix+"%s\n", v.PrettyPrint())
		case OutputDefinition:
			fmt.Printf(prefix+"%s\n", v.PrettyPrint())

		}
	}

}

//
// DEPRECATED (MOVED ELSEWHERE)
//
// type Event struct {
// 	Payload string // changeme
// //	Stamp []Tag
// }
// //
// // eval(Event e) bool
// //
// func (p AndPredicate) Eval(e Event) bool {
// 	return p.Left.Eval(e) && p.Right.Eval(e)
// }
// func (p OrPredicate) Eval(e Event) bool {
// 	return p.Left.Eval(e) || p.Right.Eval(e)
// }
// func (p NotPredicate) Eval(e Event) bool {
// 	return !p.Inner.Eval(e)
// }
// func (p TruePredicate) Eval(e Event) bool {
// 	return true
// }
// func (p FalsePredicate) Eval(e Event) bool {
// 	return false
// }

//type Monitor Filters

type Tag struct {
	//	Tag dt.Channel
	Tag string
}

type Identifier struct {
	Val string
}

type PathName struct {
	Val string
}
type QuotedString struct {
	Val string
}

type Alphanum struct {
	Val string
}

type Keyword struct {
	Val string
}

func NewIdentifier(s string) Identifier {
	return Identifier{s}
}
func NewPathName(s string) PathName {
	return PathName{s}
}
func NewQuotedString(s string) QuotedString {
	return QuotedString{s}
}

func ToSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}
