package dLola

import (
	"errors"
	"fmt"
	"strconv"
)

//Expression
type InstExpr interface {
	Sprint() string
	//	Substitute(string, interface{})
}
type InstConstExpr struct { // implements Expr,NumExpr,BoolExpr
	Name StreamName
	Pos  Position
}

type InstLetExpr struct {
	Name StreamName
	Bind InstExpr
	Body InstExpr
}
type InstIfThenElseExpr struct { // implements Expr,NumExpr,BoolExpr
	If   InstExpr
	Then InstExpr
	Else InstExpr
}
type InstStreamOffsetExpr struct { // StreamOffsetExpr implements Expr,NumExpr,BoolExpr, StrExpr
	SExpr InstStreamExpr
}
type InstBooleanExpr struct {
	BExpr InstBoolExpr
}
type InstNumericExpr struct {
	NExpr InstNumExpr
}
type InstStringExpr struct {
	StExpr InstStrExpr
}

//Stream
type InstStreamExpr interface {
	Sprint() string
}

type InstStreamFetchExpr struct { //implements StreamExpr
	Name StreamName
	Tick int
	//Default DefaultExpr //default value for the instantiated stream that gets out of the trace
	//Pos     Position
}

//Boolean
type InstBoolExpr interface {
	Sprint() string
	//	Substitute(string, BoolExpr) //TruePredicate or FalsePredicate
}

type InstTruePredicate struct{ Pos Position }
type InstFalsePredicate struct{ Pos Position }

type InstNotPredicate struct {
	Inner InstBoolExpr
}
type InstAndPredicate struct {
	Left  InstBoolExpr
	Right InstBoolExpr
}
type InstOrPredicate struct {
	Left  InstBoolExpr
	Right InstBoolExpr
}
type InstIfThenElsePredicate struct {
	If   InstBoolExpr
	Then InstBoolExpr
	Else InstBoolExpr
}
type InstNumComparisonPredicate struct {
	Comp InstNumComparison
}
type InstStrComparisonPredicate struct {
	Comp InstStrComparison
}

//Numeric
type InstNumComparison interface {
	Sprint() string
}

type InstNumLess struct {
	Left  InstNumExpr
	Right InstNumExpr
}

type InstNumLessEq struct {
	Left  InstNumExpr
	Right InstNumExpr
}

type InstNumEq struct {
	Left  InstNumExpr
	Right InstNumExpr
}

type InstNumGreater struct {
	Left  InstNumExpr
	Right InstNumExpr
}

type InstNumGreaterEq struct {
	Left  InstNumExpr
	Right InstNumExpr
}

type InstNumNotEq struct {
	Left  InstNumExpr
	Right InstNumExpr
}
type InstNumExpr interface {
	Sprint() string
}

type InstIntLiteralExpr struct {
	Num int
	//	Pos Position
}
type InstFloatLiteralExpr struct {
	Num float32
	//	Pos Position
}
type InstNumMulExpr struct {
	Left  InstNumExpr
	Right InstNumExpr
}
type InstNumDivExpr struct {
	Left  InstNumExpr
	Right InstNumExpr
}
type InstNumPlusExpr struct {
	Left  InstNumExpr
	Right InstNumExpr
}
type InstNumMinusExpr struct {
	Left  InstNumExpr
	Right InstNumExpr
}

//String
type InstStrExpr interface {
	Sprint() string
}

type InstStringLiteralExpr struct {
	S string
	//	Pos Position
}
type InstStrConcatExpr struct {
	Left  InstStrExpr
	Right InstStrExpr
}
type InstStrComparison interface {
	Sprint() string
}

type InstStrEqExpr struct {
	Left  InstStrExpr
	Right InstStrExpr
}

func (e ConstExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstConstExpr{e.Name, e.Pos}
}
func (e LetExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstLetExpr{e.Name, e.Bind.InstantiateExpr(tick, tlen), e.Body.InstantiateExpr(tick, tlen)}
}
func (e IfThenElseExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstIfThenElseExpr{e.If.InstantiateExpr(tick, tlen), e.Then.InstantiateExpr(tick, tlen), e.Else.InstantiateExpr(tick, tlen)}
}
func (e StreamOffsetExpr) InstantiateExpr(tick, tlen int) InstExpr { // expr = a[x|d] we will use the default value to infer the type
	return e.SExpr.InstantiateStreamExpr(tick, tlen) //note it does not follow the pattern of the rest
}
func (e BooleanExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstBooleanExpr{e.BExpr.InstantiateBoolExpr(tick, tlen)}
}

func (e NumericExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstNumericExpr{e.NExpr.InstantiateNumExpr(tick, tlen)}
}
func (e StringExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstStringExpr{e.StExpr.InstantiateStrExpr(tick, tlen)}
}

//Boolean
func (e TruePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstTruePredicate{}
}
func (e FalsePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstFalsePredicate{}
}
func (e NotPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstNotPredicate{e.Inner.InstantiateBoolExpr(tick, tlen)}
}
func (s StreamOffsetExpr) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstStreamOffsetExpr{s.SExpr.InstantiateBoolStreamExpr(tick, tlen)}
}
func (s ConstExpr) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstConstExpr{s.Name, s.Pos}
}
func (e AndPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstAndPredicate{e.Left.InstantiateBoolExpr(tick, tlen), e.Right.InstantiateBoolExpr(tick, tlen)}
}

func (e OrPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstOrPredicate{e.Left.InstantiateBoolExpr(tick, tlen), e.Right.InstantiateBoolExpr(tick, tlen)}
}

/*func (e IfThenElsePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstIfThenElsePredicate{e.If.InstantiateBoolExpr(tick, tlen), e.Then.InstantiateBoolExpr(tick, tlen), e.Else.InstantiateBoolExpr(tick, tlen)}
}*/
func (e NumComparisonPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstNumComparisonPredicate{e.Comp.InstantiateNumCompExpr(tick, tlen)}
}
func (e StrComparisonPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstStrComparisonPredicate{e.Comp.InstantiateStrCompExpr(tick, tlen)}
}

//Stream
func (s StreamFetchExpr) InstantiateStreamExpr(tick, tlen int) InstExpr {
	if s.Offset.val+tick < 0 || s.Offset.val+tick > tlen {
		return convertToInstExpr(s.Default)
	}
	r := InstStreamOffsetExpr{InstStreamFetchExpr{s.Name, s.Offset.val + tick}}
	//	fmt.Printf("Instantiated stream: %s for tick %d with tlen %d\n", r.Sprint(), tick, tlen)
	return r
}
func convertToInstExpr(d DefaultExpr) InstExpr {
	var r InstExpr
	switch v := d.val.(type) {
	case TruePredicate:
		r = InstTruePredicate{}
	case FalsePredicate:
		r = InstFalsePredicate{}
	case IntLiteralExpr:
		r = InstIntLiteralExpr{v.Num}
	case FloatLiteralExpr:
		r = InstFloatLiteralExpr{v.Num}
	case StringLiteralExpr:
		r = InstStringLiteralExpr{v.S}
	default: //will occurr for streams used as s = r without specifying default values
		errors.New("Impossible to convert other thing to InstExpr\n")
	}
	return r

}

func (s StreamFetchExpr) InstantiateBoolStreamExpr(tick, tlen int) InstBoolExpr {
	if s.Offset.val+tick < 0 || s.Offset.val+tick > tlen {
		return convertToInstExpr(s.Default).(InstBoolExpr)
	}
	r := InstStreamOffsetExpr{InstStreamFetchExpr{s.Name, s.Offset.val + tick}}
	//	fmt.Printf("Instantiated stream: %s for tick %d with tlen %d\n", r.Sprint(), tick, tlen)
	return r

}
func (s StreamFetchExpr) InstantiateNumStreamExpr(tick, tlen int) InstNumExpr {
	if s.Offset.val+tick < 0 || s.Offset.val+tick > tlen {
		return convertToInstExpr(s.Default).(InstNumExpr)
	}
	r := InstStreamOffsetExpr{InstStreamFetchExpr{s.Name, s.Offset.val + tick}}
	//	fmt.Printf("Instantiated stream: %s for tick %d with tlen %d\n", r.Sprint(), tick, tlen)
	return r

}
func (s StreamFetchExpr) InstantiateStrStreamExpr(tick, tlen int) InstStrExpr {
	if s.Offset.val+tick < 0 || s.Offset.val+tick > tlen {
		return convertToInstExpr(s.Default).(InstStrExpr)
	}
	r := InstStreamOffsetExpr{InstStreamFetchExpr{s.Name, s.Offset.val + tick}}
	//	fmt.Printf("Instantiated stream: %s for tick %d with tlen %d\n", r.Sprint(), tick, tlen)
	return r
}

//Num
func (this NumLess) InstantiateNumCompExpr(tick, tlen int) InstNumComparison {
	return InstNumLess{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (this NumLessEq) InstantiateNumCompExpr(tick, tlen int) InstNumComparison {
	return InstNumLessEq{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (this NumGreater) InstantiateNumCompExpr(tick, tlen int) InstNumComparison {
	return InstNumGreater{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (this NumGreaterEq) InstantiateNumCompExpr(tick, tlen int) InstNumComparison {
	return InstNumGreaterEq{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (this NumEq) InstantiateNumCompExpr(tick, tlen int) InstNumComparison {
	return InstNumEq{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (this NumNotEq) InstantiateNumCompExpr(tick, tlen int) InstNumComparison {
	return InstNumNotEq{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}

func (this IntLiteralExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstIntLiteralExpr{this.Num}
}
func (this FloatLiteralExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstFloatLiteralExpr{this.Num}
}
func (this NumMulExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstNumMulExpr{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (this NumDivExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstNumDivExpr{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (this NumPlusExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstNumPlusExpr{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (this NumMinusExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstNumMinusExpr{this.Left.InstantiateNumExpr(tick, tlen), this.Right.InstantiateNumExpr(tick, tlen)}
}
func (s StreamOffsetExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstStreamOffsetExpr{s.SExpr.InstantiateNumStreamExpr(tick, tlen)}
}
func (s ConstExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstConstExpr{s.Name, s.Pos}
}

//String
func (this StringLiteralExpr) InstantiateStrExpr(tick, tlen int) InstStrExpr {
	return InstStringLiteralExpr{this.S}
}
func (this StrConcatExpr) InstantiateStrExpr(tick, tlen int) InstStrExpr {
	return InstStrConcatExpr{this.Left.InstantiateStrExpr(tick, tlen), this.Right.InstantiateStrExpr(tick, tlen)}
}
func (this StreamOffsetExpr) InstantiateStrExpr(tick, tlen int) InstStrExpr {
	return InstStreamOffsetExpr{this.SExpr.InstantiateStrStreamExpr(tick, tlen)}
}
func (this ConstExpr) InstantiateStrExpr(tick, tlen int) InstStrExpr {
	return InstConstExpr{this.Name, this.Pos}
}

func (this StrEqExpr) InstantiateStrCompExpr(tick, tlen int) InstStrComparison {
	return InstStrEqExpr{this.Left.InstantiateStrExpr(tick, tlen), this.Right.InstantiateStrExpr(tick, tlen)}
}

//Sprint
//Expr
func (this InstConstExpr) Sprint() string {
	return string(this.Name)
}

func (this InstLetExpr) Sprint() string {
	bind := this.Bind.Sprint()
	body := this.Bind.Sprint()
	return fmt.Sprintf("let %s = %s in %s", this.Name, bind, body)
}
func (this InstIfThenElseExpr) Sprint() string {
	if_part := this.If.Sprint()
	then_part := this.Then.Sprint()
	else_part := this.Else.Sprint()
	return fmt.Sprintf("if %s then %s else %s", if_part, then_part, else_part)
}
func (this InstNumericExpr) Sprint() string {
	return this.NExpr.Sprint()
}
func (this InstStreamOffsetExpr) Sprint() string {
	return this.SExpr.Sprint()
}
func (this InstBooleanExpr) Sprint() string {
	return this.BExpr.Sprint()
}
func (this InstStringExpr) Sprint() string {
	return this.StExpr.Sprint()
}

//Stream
func (e InstStreamFetchExpr) Sprint() string {
	return fmt.Sprintf("%s[%d]", e.Name.Sprint(), e.Tick)
}

//predicates
func (p InstAndPredicate) Sprint() string {
	return fmt.Sprintf("(%s) /\\ (%s)", p.Left.Sprint(), p.Right.Sprint())
}
func (p InstOrPredicate) Sprint() string {
	return fmt.Sprintf("(%s) \\/  (%s)", p.Left.Sprint(), p.Right.Sprint())
}
func (p InstNotPredicate) Sprint() string {
	return fmt.Sprintf("~ (%s)", p.Inner.Sprint())
}
func (p InstTruePredicate) Sprint() string {
	return fmt.Sprintf("true")
}
func (p InstFalsePredicate) Sprint() string {
	return fmt.Sprintf("false")
}
func (p InstNumComparisonPredicate) Sprint() string {
	return p.Comp.Sprint()
}
func (p InstStrComparisonPredicate) Sprint() string {
	return p.Comp.Sprint()
}

//numeric
func (this InstNumLess) Sprint() string {
	return fmt.Sprintf("(%s) < (%s)", this.Left.Sprint(), this.Right.Sprint())
}
func (this InstNumLessEq) Sprint() string {
	return fmt.Sprintf("(%s) <= (%s)", this.Left.Sprint(), this.Right.Sprint())
}
func (this InstNumGreater) Sprint() string {
	return fmt.Sprintf("(%s) > (%s)", this.Left.Sprint(), this.Right.Sprint())
}
func (this InstNumGreaterEq) Sprint() string {
	return fmt.Sprintf("(%s) >= (%s)", this.Left.Sprint(), this.Right.Sprint())
}
func (this InstNumNotEq) Sprint() string {
	return fmt.Sprintf("(%s) != (%s)", this.Left.Sprint(), this.Right.Sprint())
}
func (this InstNumEq) Sprint() string {
	return fmt.Sprintf("(%s) = (%s)", this.Left.Sprint(), this.Right.Sprint())
}

func (e InstNumMulExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)", e.Left.Sprint(), "*", e.Right.Sprint())
}
func (e InstNumDivExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)", e.Left.Sprint(), "/", e.Right.Sprint())
}
func (e InstNumPlusExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)", e.Left.Sprint(), "+", e.Right.Sprint())
}
func (e InstNumMinusExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)", e.Left.Sprint(), "-", e.Right.Sprint())
}
func (e InstIntLiteralExpr) Sprint() string {
	return strconv.Itoa(e.Num)
}
func (e InstFloatLiteralExpr) Sprint() string {
	return strconv.FormatFloat(float64(e.Num), 'f', 4, 32)
}

//string
func (this InstStringLiteralExpr) Sprint() string {
	return this.S
}
func (this InstStrConcatExpr) Sprint() string {
	return fmt.Sprintf("(%s) strConcat (%s)", this.Left.Sprint(), this.Right.Sprint())
}

func (this InstStrEqExpr) Sprint() string {
	return fmt.Sprintf("(%s) strEq (%s)", this.Left.Sprint(), this.Right.Sprint())
}
