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
	LastType = Unknown
)

type StreamName string

func (s StreamName) Sprint() string {
	return string(s)
}

func (t StreamType) Sprint() string {
	type_names := []string{"num", "bool", "string"}
	if t >= LastType {
		return ""
	}
	return fmt.Sprintf("%s", type_names[t])
}

type Position struct {
	Line, Col, Offset int
}

type ConstDecl struct { // const int one_sec := 1s
	Name StreamName
	Type StreamType
	Val  Expr
	Pos  Position
}
type InputDecl struct { // input int bar
	Name StreamName
	Type StreamType
	Pos  Position
}
type OutputDecl struct { // output int foo /* this is just a decl, later a tick and a def will be given */
	Name StreamName
	Type StreamType
	Pos  Position
}

type MonitorDecl struct {
	Nid   int
	Decls []interface{}
}

func NewMonitorDecl(n, d interface{}) MonitorDecl {
	return MonitorDecl{n.(IntLiteralExpr).Num, ToSlice(d)}
}

type TopoMonitorDecls struct {
	Topo         string
	Nmons        int
	MonitorDecls []MonitorDecl
}

func NewTopoMonitorDecls(t, m interface{}) TopoMonitorDecls {
	monitorDecls := make([]MonitorDecl, 0)
	nmons := 0
	for _, m := range ToSlice(m) {
		monitorDecls = append(monitorDecls, m.(MonitorDecl))
		nmons++
	}
	return TopoMonitorDecls{t.(Identifier).Val, nmons, monitorDecls}
}

/*lm: not needed type TicksDecl struct {
	Name  StreamName
	Ticks TickingExpr
}*/

type OutputDefinition struct {
	Name StreamName
	Type StreamType
	Eval bool
	Expr Expr // chango to ValueExpr?
	Pos  Position
}

func NewPosition(p interface{}) Position {
	po := p.(position) //this type is defined in parser.go as part of the PIGEON library
	return Position{po.line, po.col, po.offset}
}

func NewConstDecl(n, t, e, p interface{}) ConstDecl {
	name := getStreamName(n)
	return ConstDecl{name, t.(StreamType), e.(Expr), NewPosition(p)}
}
func NewInputDecl(n, t, p interface{}) InputDecl {
	name := getStreamName(n)
	return InputDecl{name, t.(StreamType), NewPosition(p)}
}
func NewOutputDecl(n, t, p interface{}) OutputDecl {
	name := getStreamName(n)
	return OutputDecl{name, t.(StreamType), NewPosition(p)}
}

/*: not needed func NewTicksDecl(n, t interface{}) TicksDecl {
	name := getStreamName(n)
	expr := t.(TickingExpr)
	return TicksDecl{name, expr}
}*/

func NewOutputDefinition(n, t, le, e, p interface{}) OutputDefinition {
	name := getStreamName(n)
	expr := e.(Expr)
	eval := getEval(le)
	return OutputDefinition{name, t.(StreamType), eval, expr, NewPosition(p)}
}

func getStreamName(a interface{}) StreamName {
	return StreamName(a.(Identifier).Val)
}

func getEval(le interface{}) bool {
	eval := true
	v, ok := le.(bool)
	if ok {
		eval = v
	} //if le was not a bool, it will be considered eval
	return eval
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
