package dLola

import (
	"errors"
	"fmt"
	"strconv"
)

//Expression
type InstExpr interface {
	Sprint() string
	Substitute(InstStreamExpr, InstExpr) InstExpr
	Simplify() (InstExpr, bool)
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
	Substitute(InstStreamExpr, InstExpr) InstExpr
	SubstituteBool(InstStreamExpr, InstExpr) InstBoolExpr
	SubstituteNum(InstStreamExpr, InstExpr) InstNumExpr
	SubstituteStr(InstStreamExpr, InstExpr) InstStrExpr
	GetName() StreamName
	GetTick() int
	//Simplify() InstExpr THESE WILL NOT BE USED AS AN INSTANTIATED STREAM CANNOT BE SIMPLIFIED
	//SimplifyBool() InstBoolExpr
	//SimplifyNum() InstNumExpr
	//SimplifyStr() InstStrExpr
}
type InstStreamFetchExpr struct { //implements StreamExpr
	Name StreamName
	Tick int
	//Default DefaultExpr //default value for the instantiated stream that gets out of the trace
	//Pos     Position
}

func (this InstStreamFetchExpr) GetName() StreamName {
	return this.Name
}
func (this InstStreamFetchExpr) GetTick() int {
	return this.Tick
}

//Boolean
type InstBoolExpr interface {
	Sprint() string
	SubstituteBool(InstStreamExpr, InstExpr) InstBoolExpr
	SimplifyBool() (InstBoolExpr, bool)
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
	SubstituteNumComp(InstStreamExpr, InstExpr) InstNumComparison
	SimplifyNumComp() (InstBoolExpr, bool)
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
	SubstituteNum(InstStreamExpr, InstExpr) InstNumExpr
	SimplifyNum() (InstNumExpr, bool)
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
	SubstituteStr(InstStreamExpr, InstExpr) InstStrExpr
	SimplifyStr() (InstStrExpr, bool)
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
	SubstituteStrComp(InstStreamExpr, InstExpr) InstStrComparison
	SimplifyStrComp() (InstBoolExpr, bool)
}
type InstStrEqExpr struct {
	Left  InstStrExpr
	Right InstStrExpr
}

func (this ConstExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstConstExpr{this.Name, this.Pos}
}
func (this LetExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstLetExpr{this.Name, this.Bind.InstantiateExpr(tick, tlen), this.Body.InstantiateExpr(tick, tlen)}
}
func (this IfThenElseExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstIfThenElseExpr{this.If.InstantiateExpr(tick, tlen), this.Then.InstantiateExpr(tick, tlen), this.Else.InstantiateExpr(tick, tlen)}
}
func (this StreamOffsetExpr) InstantiateExpr(tick, tlen int) InstExpr { // expr = a[x|d] we will use the default value to infer the type
	return this.SExpr.InstantiateStreamExpr(tick, tlen) //note it does not follow the pattern of the rest
}
func (this BooleanExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstBooleanExpr{this.BExpr.InstantiateBoolExpr(tick, tlen)}
}
func (this NumericExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstNumericExpr{this.NExpr.InstantiateNumExpr(tick, tlen)}
}
func (this StringExpr) InstantiateExpr(tick, tlen int) InstExpr {
	return InstStringExpr{this.StExpr.InstantiateStrExpr(tick, tlen)}
}

//Boolean
func (this TruePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstTruePredicate{}
}
func (this FalsePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstFalsePredicate{}
}
func (this NotPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstNotPredicate{this.Inner.InstantiateBoolExpr(tick, tlen)}
}
func (this StreamOffsetExpr) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return this.SExpr.InstantiateBoolStreamExpr(tick, tlen)
}
func (this ConstExpr) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstConstExpr{this.Name, this.Pos}
}
func (this AndPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstAndPredicate{this.Left.InstantiateBoolExpr(tick, tlen), this.Right.InstantiateBoolExpr(tick, tlen)}
}
func (this OrPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstOrPredicate{this.Left.InstantiateBoolExpr(tick, tlen), this.Right.InstantiateBoolExpr(tick, tlen)}
}

/*func (this IfThenElsePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstIfThenElsePredicate{this.If.InstantiateBoolExpr(tick, tlen), this.Then.InstantiateBoolExpr(tick, tlen), this.Else.InstantiateBoolExpr(tick, tlen)}
}*/
func (this NumComparisonPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstNumComparisonPredicate{this.Comp.InstantiateNumCompExpr(tick, tlen)}
}
func (this StrComparisonPredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstStrComparisonPredicate{this.Comp.InstantiateStrCompExpr(tick, tlen)}
}

//Stream
func (this StreamFetchExpr) InstantiateStreamExpr(tick, tlen int) InstExpr {
	if this.Offset.val+tick < 0 || this.Offset.val+tick > tlen {
		return convertToInstExpr(this.Default)
	}
	r := InstStreamOffsetExpr{InstStreamFetchExpr{this.Name, this.Offset.val + tick}}
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
func (this StreamFetchExpr) InstantiateBoolStreamExpr(tick, tlen int) InstBoolExpr {
	if this.Offset.val+tick < 0 || this.Offset.val+tick > tlen {
		return convertToInstExpr(this.Default).(InstBoolExpr)
	}
	r := InstStreamOffsetExpr{InstStreamFetchExpr{this.Name, this.Offset.val + tick}}
	//	fmt.Printf("Instantiated stream: %s for tick %d with tlen %d\n", r.Sprint(), tick, tlen)
	return r
}
func (this StreamFetchExpr) InstantiateNumStreamExpr(tick, tlen int) InstNumExpr {
	if this.Offset.val+tick < 0 || this.Offset.val+tick > tlen {
		return convertToInstExpr(this.Default).(InstNumExpr)
	}
	r := InstStreamOffsetExpr{InstStreamFetchExpr{this.Name, this.Offset.val + tick}}
	//	fmt.Printf("Instantiated stream: %s for tick %d with tlen %d\n", r.Sprint(), tick, tlen)
	return r
}
func (this StreamFetchExpr) InstantiateStrStreamExpr(tick, tlen int) InstStrExpr {
	if this.Offset.val+tick < 0 || this.Offset.val+tick > tlen {
		return convertToInstExpr(this.Default).(InstStrExpr)
	}
	r := InstStreamOffsetExpr{InstStreamFetchExpr{this.Name, this.Offset.val + tick}}
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
func (this StreamOffsetExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return this.SExpr.InstantiateNumStreamExpr(tick, tlen)
}
func (this ConstExpr) InstantiateNumExpr(tick, tlen int) InstNumExpr {
	return InstConstExpr{this.Name, this.Pos}
}

//String
func (this StringLiteralExpr) InstantiateStrExpr(tick, tlen int) InstStrExpr {
	return InstStringLiteralExpr{this.S}
}
func (this StrConcatExpr) InstantiateStrExpr(tick, tlen int) InstStrExpr {
	return InstStrConcatExpr{this.Left.InstantiateStrExpr(tick, tlen), this.Right.InstantiateStrExpr(tick, tlen)}
}
func (this StreamOffsetExpr) InstantiateStrExpr(tick, tlen int) InstStrExpr {
	return this.SExpr.InstantiateStrStreamExpr(tick, tlen)
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
func (this InstStreamFetchExpr) Sprint() string {
	return fmt.Sprintf("%s[%d]", this.Name.Sprint(), this.Tick)
}

//predicates
func (this InstAndPredicate) Sprint() string {
	return fmt.Sprintf("(%s) /\\ (%s)", this.Left.Sprint(), this.Right.Sprint())
}
func (this InstOrPredicate) Sprint() string {
	return fmt.Sprintf("(%s) \\/  (%s)", this.Left.Sprint(), this.Right.Sprint())
}
func (this InstNotPredicate) Sprint() string {
	return fmt.Sprintf("~ (%s)", this.Inner.Sprint())
}
func (this InstTruePredicate) Sprint() string {
	return fmt.Sprintf("true")
}
func (this InstFalsePredicate) Sprint() string {
	return fmt.Sprintf("false")
}
func (this InstNumComparisonPredicate) Sprint() string {
	return this.Comp.Sprint()
}
func (this InstStrComparisonPredicate) Sprint() string {
	return this.Comp.Sprint()
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

func (this InstNumMulExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)", this.Left.Sprint(), "*", this.Right.Sprint())
}
func (this InstNumDivExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)", this.Left.Sprint(), "/", this.Right.Sprint())
}
func (this InstNumPlusExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)", this.Left.Sprint(), "+", this.Right.Sprint())
}
func (this InstNumMinusExpr) Sprint() string {
	return fmt.Sprintf("(%s)%s(%s)", this.Left.Sprint(), "-", this.Right.Sprint())
}
func (this InstIntLiteralExpr) Sprint() string {
	return strconv.Itoa(this.Num)
}
func (this InstFloatLiteralExpr) Sprint() string {
	return strconv.FormatFloat(float64(this.Num), 'f', 4, 32)
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

//Substitute
func (this InstConstExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return this
}
func (this InstLetExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return InstLetExpr{this.Name, this.Bind.Substitute(s, v), this.Body.Substitute(s, v)}
}
func (this InstIfThenElseExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return InstIfThenElseExpr{this.If.Substitute(s, v), this.Then.Substitute(s, v), this.Else.Substitute(s, v)}
}
func (this InstStreamOffsetExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return this.SExpr.Substitute(s, v) //note it does not follow the pattern of the rest
}
func (this InstBooleanExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return InstBooleanExpr{this.BExpr.SubstituteBool(s, v)}
}
func (this InstNumericExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return InstNumericExpr{this.NExpr.SubstituteNum(s, v)}
}
func (this InstStringExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return InstStringExpr{this.StExpr.SubstituteStr(s, v)}
}

//Boolean
func (this InstTruePredicate) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return this
}
func (this InstFalsePredicate) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return this
}
func (this InstNotPredicate) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return InstNotPredicate{this.Inner.SubstituteBool(s, v)}
}
func (this InstStreamOffsetExpr) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return this.SExpr.SubstituteBool(s, v)
}
func (this InstConstExpr) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return this
}
func (this InstAndPredicate) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return InstAndPredicate{this.Left.SubstituteBool(s, v), this.Right.SubstituteBool(s, v)}
}
func (this InstOrPredicate) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return InstOrPredicate{this.Left.SubstituteBool(s, v), this.Right.SubstituteBool(s, v)}
}

/*func (this InstIfThenElsePredicate) InstantiateBoolExpr(tick, tlen int) InstBoolExpr {
	return InstIfThenElsePredicate{this.If.Substitute(s, v), this.Then.Substitute(s, v), this.Else.Substitute(s, v)}
}*/
func (this InstNumComparisonPredicate) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return InstNumComparisonPredicate{this.Comp.SubstituteNumComp(s, v)}
}
func (this InstStrComparisonPredicate) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	return InstStrComparisonPredicate{this.Comp.SubstituteStrComp(s, v)}
}

//Stream
func (this InstStreamFetchExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	if s.GetName() == this.Name && this.Tick == s.GetTick() {
		return v
	}
	return InstStreamOffsetExpr{this}
}
func (this InstStreamFetchExpr) SubstituteBool(s InstStreamExpr, v InstExpr) InstBoolExpr {
	if s.GetName() == this.Name && this.Tick == s.GetTick() {
		return v.(InstBoolExpr)
	}
	return InstStreamOffsetExpr{this}
}
func (this InstStreamFetchExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	if s.GetName() == this.Name && this.Tick == s.GetTick() {
		return v.(InstNumExpr)
	}
	return InstStreamOffsetExpr{this}
}
func (this InstStreamFetchExpr) SubstituteStr(s InstStreamExpr, v InstExpr) InstStrExpr {
	if s.GetName() == this.Name && this.Tick == s.GetTick() {
		return v.(InstStrExpr)
	}
	return InstStreamOffsetExpr{this}
}

//Num
func (this InstNumLess) SubstituteNumComp(s InstStreamExpr, v InstExpr) InstNumComparison {
	return InstNumLess{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstNumLessEq) SubstituteNumComp(s InstStreamExpr, v InstExpr) InstNumComparison {
	return InstNumLessEq{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstNumGreater) SubstituteNumComp(s InstStreamExpr, v InstExpr) InstNumComparison {
	return InstNumGreater{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstNumGreaterEq) SubstituteNumComp(s InstStreamExpr, v InstExpr) InstNumComparison {
	return InstNumGreaterEq{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstNumEq) SubstituteNumComp(s InstStreamExpr, v InstExpr) InstNumComparison {
	return InstNumEq{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstNumNotEq) SubstituteNumComp(s InstStreamExpr, v InstExpr) InstNumComparison {
	return InstNumNotEq{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}

func (this InstIntLiteralExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	return this
}
func (this InstFloatLiteralExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	return this
}
func (this InstNumMulExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	return InstNumMulExpr{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstNumDivExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	return InstNumDivExpr{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstNumPlusExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	return InstNumPlusExpr{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstNumMinusExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	return InstNumMinusExpr{this.Left.SubstituteNum(s, v), this.Right.SubstituteNum(s, v)}
}
func (this InstStreamOffsetExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	return this.SExpr.SubstituteNum(s, v)
}
func (this InstConstExpr) SubstituteNum(s InstStreamExpr, v InstExpr) InstNumExpr {
	return this
}

//String
func (this InstStringLiteralExpr) SubstituteStr(s InstStreamExpr, v InstExpr) InstStrExpr {
	return this
}
func (this InstStrConcatExpr) SubstituteStr(s InstStreamExpr, v InstExpr) InstStrExpr {
	return InstStrConcatExpr{this.Left.SubstituteStr(s, v), this.Right.SubstituteStr(s, v)}
}
func (this InstStreamOffsetExpr) SubstituteStr(s InstStreamExpr, v InstExpr) InstStrExpr {
	return this.SExpr.SubstituteStr(s, v)
}
func (this InstConstExpr) SubstituteStr(s InstStreamExpr, v InstExpr) InstStrExpr {
	return this
}
func (this InstStrEqExpr) SubstituteStrComp(s InstStreamExpr, v InstExpr) InstStrComparison {
	return InstStrEqExpr{this.Left.SubstituteStr(s, v), this.Right.SubstituteStr(s, v)}
}

//Literals need to implement InstExpr to compile, implementation of Substitute (should not be needed at runtime)
func (this InstTruePredicate) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return this
}
func (this InstFalsePredicate) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return this
}
func (this InstIntLiteralExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return this
}
func (this InstFloatLiteralExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return this
}
func (this InstStringLiteralExpr) Substitute(s InstStreamExpr, v InstExpr) InstExpr {
	return this
}
