package dLola

import (
	"errors"
	"fmt"
)

// parsed Output
type OutputStream struct {
	Name StreamName
	Type StreamType
	//	Ticks TickingExpr
	Expr Expr // chango to ValueExpr?
}

// TODO
// type Symbol {
//	Name StreamName
//	Type StreamType
//	Expr * the_expr
// }
//
// var SymbolTable [StreamName]Symbol

type Spec struct {
	Input  map[StreamName]InputDecl
	Const  map[StreamName]ConstDecl
	Output map[StreamName]OutputStream
}

func newSpec() *Spec {
	spec := Spec{}
	spec.Input = make(map[StreamName]InputDecl)
	spec.Const = make(map[StreamName]ConstDecl)
	spec.Output = make(map[StreamName]OutputStream)
	return &spec
}

//
// after ParseFile returns a []interafce{} all the elements in the slice
//    are Filer or Session or Stream or Trigger
//    this function creates a MonitorMachine from such a mixed slice
//
type specInProgress struct {
	Output map[StreamName]OutputDecl
	//lm: not needed	Ticks  map[StreamName]TicksDecl
	Define map[StreamName]OutputDefinition
}

func newSpecInProgress() *specInProgress {
	s := specInProgress{}
	s.Output = make(map[StreamName]OutputDecl)
	//	s.Ticks = make(map[StreamName]TicksDecl)
	s.Define = make(map[StreamName]OutputDefinition)

	return &s
}

func declared_any(name StreamName, spec *Spec, prog *specInProgress) error {
	_, present_const := spec.Const[name]
	_, present_input := spec.Input[name]
	_, present_output := prog.Output[name]
	//	_, present_ticks := prog.Ticks[name]
	_, present_define := prog.Define[name]
	if present_const || present_input || present_output /*|| present_ticks*/ || present_define {
		str := fmt.Sprintf("%s already declared", string(name))
		return errors.New(str)
	}
	return nil
}

func declared_input(name StreamName, spec *Spec) bool {
	_, present := spec.Input[name]
	return present
}
func declared_const(name StreamName, spec *Spec) bool {
	_, present := spec.Const[name]
	return present
}
func declared_output(name StreamName, spec *specInProgress) bool {
	_, present := spec.Output[name]
	return present
}

/*lm: not needed func declared_ticks(name StreamName, spec *specInProgress) bool {
	_, present := spec.Ticks[name]
	return present
}*/
func declared_define(name StreamName, spec *specInProgress) bool {
	_, present := spec.Define[name]
	return present
}

func ProcessDeclarations(ds []interface{}) (*Spec, error) {

	spec := newSpec()
	in_progress := newSpecInProgress()

	for _, v := range ds {
		switch decl := v.(type) {
		case InputDecl:
			name := StreamName(decl.Name)
			if err := declared_any(name, spec, in_progress); err != nil {
				return nil, err
			}
			spec.Input[name] = decl
		case ConstDecl:
			name := StreamName(decl.Name)
			if err := declared_any(name, spec, in_progress); err != nil {
				return nil, err
			}
			spec.Const[name] = decl
		case OutputDecl:
			name := StreamName(decl.Name)
			if declared_input(name, spec) ||
				declared_const(name, spec) ||
				declared_output(name, in_progress) {
				str := fmt.Sprintf("%s redeclared", name)
				return nil, errors.New(str)
			}
			in_progress.Output[name] = decl
		/*case TicksDecl:
		name := StreamName(decl.Name)
		if declared_input(name, spec) ||
			declared_const(name, spec) ||
			declared_ticks(name, in_progress) {
			str := fmt.Sprintf("%s redeclared", name)
			return nil, errors.New(str)
		}
		in_progress.Ticks[name] = decl
		*/
		case OutputDefinition:
			name := StreamName(decl.Name)
			if declared_input(name, spec) ||
				declared_const(name, spec) ||
				declared_define(name, in_progress) {
				str := fmt.Sprintf("%s redeclared", name)
				return nil, errors.New(str)
			}
			in_progress.Output[name] = OutputDecl{name, decl.Type, decl.Pos} //the sentence output num a = 2 will combine output declaration and definition, it is intrinsically declared by this line
			in_progress.Define[name] = decl
		case string:
			//ignore it is a comment
		default:
			str := fmt.Sprint("Unexpected type returned by parser: %t", v)
			return nil, errors.New(str)
		}
	}
	//
	//  1.Check that all output streams appear in ticks and defined
	//  exactly once
	for key, decl := range in_progress.Output {
		//tick, is_tick := in_progress.Ticks[key]
		def, is_define := in_progress.Define[key]
		/*if !is_tick { // "output" but not "ticks"
			str := fmt.Sprintf("stream %s is defined as\"output\" but not \"ticks\"\n", key)
			return spec, errors.New(str)
		}*/
		if !is_define { // "output" but not "define"
			str := fmt.Sprintf("stream %s is defined as\"output\" but not \"define\"\n", key)
			return spec, errors.New(str)
			if def.Type != decl.Type { // inconsistent types
				str := fmt.Sprintf("%s has diferent types in \"output\" and \"define\": %s and %s\n", key, decl.Type.Sprint(), def.Type.Sprint())
				return spec, errors.New(str)
			}
		}
		// OK. All matches
		spec.Output[key] = OutputStream{key, def.Type /*, tick.Ticks*/, def.Expr}
	}

	//
	// 2. Check wether all "ticks" have "define" and "output"
	//
	/*for key, _ := range in_progress.Ticks {
		_, declared := in_progress.Output[key]
		_, defined := in_progress.Define[key]
		if !declared && !defined {
			str := fmt.Sprintf("%s has \"ticks\" but neither \"output\" nor \"define\"", key)
			return spec, errors.New(str)
		}
		if !declared {
			str := fmt.Sprintf("%s has \"ticks\" and \"define\"but not \"output\"", key)
			return spec, errors.New(str)
		}
		if !defined {
			str := fmt.Sprintf("%s has \"ticks\" and \"output\"but not \"define\"", key)
			return spec, errors.New(str)
		}
	}*/
	//
	// 3. Check wether all "define" have "ticks" and "output"
	//
	for key, _ := range in_progress.Define {
		_, declared := in_progress.Output[key]
		//_, areticks := in_progress.Ticks[key]
		/*if !declared && !areticks {
			str := fmt.Sprintf("%s has \"define\" but neither \"output\" nor \"ticks\"", key)
			return spec, errors.New(str)
		}*/
		if !declared {
			str := fmt.Sprintf("%s has \"define\" and \"ticks\"but not \"output\"", key)
			return spec, errors.New(str)
		}
		/*if !areticks {
			str := fmt.Sprintf("%s has \"define\" and \"output\"but not \"ticks\"", key)
			return spec, errors.New(str)
		}*/
	}
	return spec, nil
}

/*Auxiliary functions and methods for all the functions that operate on a Spec*/
func (c ConstDecl) Sprint() string {
	return fmt.Sprintf("ConstDecl: {Name = %s, Type = %s, expr = %s}\n", c.Name.Sprint(), c.Type.Sprint(), c.Val.Sprint())
}

func (i InputDecl) Sprint() string {
	return fmt.Sprintf("InputDecl: {Name = %s, Type = %s}\n", i.Name.Sprint(), i.Type.Sprint())
}

func (o OutputDecl) Sprint() string {
	return fmt.Sprintf("OutputDecl: {Name = %s, Type = %s}\n", o.Name.Sprint(), o.Type.Sprint())
}

func (o OutputDefinition) Sprint() string {
	return fmt.Sprintf("OutputDefinition: {Name = %s, Type = %s, expr = %s}\n", o.Name.Sprint(), o.Type.Sprint(), o.Expr.Sprint())
}

/*Pretty Print using method Accept*/
func (c ConstDecl) PrettyPrint() string {
	v := PrettyPrinterVisitor{0, ""}
	c.Val.Accept(&v)
	return fmt.Sprintf("ConstDecl: {Name = %s, Type = %s, expr = %s}\n", c.Name.Sprint(), c.Type.Sprint(), v.s)
}

func (i InputDecl) PrettyPrint() string {
	return fmt.Sprintf("InputDecl: {Name = %s, Type = %s}\n", i.Name.Sprint(), i.Type.Sprint())
}

func (o OutputDecl) PrettyPrint() string {
	return fmt.Sprintf("OutputDecl: {Name = %s, Type = %s}\n", o.Name.Sprint(), o.Type.Sprint())
}

func (o OutputDefinition) PrettyPrint() string {
	v := PrettyPrinterVisitor{0, ""}
	o.Expr.Accept(&v)
	return fmt.Sprintf("OutputDefinition: {Name = %s, Type = %s, expr = \n%s}\n", o.Name.Sprint(), o.Type.Sprint(), v.s)
}

func (o OutputStream) PrettyPrint() string {
	v := PrettyPrinterVisitor{0, ""}
	o.Expr.Accept(&v)
	return fmt.Sprintf("OutputStream: {Name = %s, Type = %s, expr = \n%s}\n", o.Name.Sprint(), o.Type.Sprint(), v.s)
}

/*Calls to the parser generated by PIGEON and returns a list of all the trees matched (each sentence will have a tree)
after that, ProcessDeclarations is called which returns a Spec */
func GetSpec(filename, prefix string) (*Spec, error) {
	ast, err := ParseFile(filename)
	if err != nil {
		fmt.Printf(prefix+"There was an error: %s\n", err)
		return ProcessDeclarations([]interface{}{})
	}
	last := ast.([]interface{})
	/*TODO: perform castings before returning -> it wont work since the output type will still be []interface{}, so direct call to methods won't work*/
	return ProcessDeclarations(last)
}

var Verbose bool = false

func Sprint(spec Spec) string {
	var str string
	if Verbose {
		str = str + fmt.Sprintf("There are %d constants\n", len(spec.Const))
		str = str + fmt.Sprintf("There are %d inputs\n", len(spec.Input))
		str = str + fmt.Sprintf("There are %d output streams\n", len(spec.Output))
	}
	for _, v := range spec.Const {
		str = str + fmt.Sprintf("const %s %s := %s\n", v.Type.Sprint(), v.Name, v.Val.Sprint())
	}
	for _, v := range spec.Input {
		str = str + fmt.Sprintf("input %s %s\n", v.Type.Sprint(), v.Name)
	}
	for _, v := range spec.Output {
		str = str + fmt.Sprintf("output %s %s\n", v.Type.Sprint(), v.Name)
		//lm: not needed str = str + fmt.Sprintf("ticks %s := %s\n", v.Name, v.Ticks.Sprint())
		str = str + fmt.Sprintf("define %s %s := %s\n", v.Type.Sprint(), v.Name, v.Expr.Sprint())
	}
	return str

}

func PrintSpec(spec *Spec, prefix string) {
	fmt.Printf(prefix + Sprint(*spec))
}

func PrettyPrintSpec(spec *Spec, prefix string) string {
	var str string
	for _, v := range spec.Const {
		str = str + v.PrettyPrint()
	}
	for _, v := range spec.Input {
		str = str + v.PrettyPrint()
	}
	for _, v := range spec.Output {
		//str = str + fmt.Sprintf("output %s %s\n", v.Type.Sprint(), v.Name)
		//lm: not needed str = str + fmt.Sprintf("ticks %s := %s\n", v.Name, v.Ticks.Sprint())
		str = str + v.PrettyPrint()
		//fmt.Sprintf("define %s %s := %s\n", v.Type.Sprint(), v.Name, v.Expr.PrettyPrint())
	}
	return str
}

func CheckTypesSpec(spec *Spec, prefix string) {
	typeVisitor := TypeVisitor{make(map[StreamName]StreamType), make([]string, 0), Unknown, false}
	for _, v := range spec.Const { //we do this first because the order in which constants are iterated over is not always the same (using the same file)
		typeVisitor.symTab[v.Name] = v.Type //introduce constant names as declared for the TypeVisitor
	}
	//constants may reference other constants
	//TODO: referencias a constantes se parsean como streams, cambiar el tipo de datos para que sea adecuado antes del checking de tipos!!!!
	//mark the typeVisitor to detect references to other streams in order to raise an error (this is a constant)
	typeVisitor.streamsForbidden = true
	for _, v := range spec.Const {
		typeVisitor.reqType = v.Type //we mark the type that the overall expression must have, the declared type of the output stream
		v.Val.Accept(&typeVisitor)
	}
	//streams are now alowed
	typeVisitor.streamsForbidden = false

	for _, v := range spec.Input {
		//introduce all the input streams as declared so when they are used TypeVisitor can know if they were declared
		typeVisitor.symTab[v.Name] = v.Type
	}

	for _, v := range spec.Output { //we do this first because the order in which output streams are iterated over is not always the same (using the same file)
		typeVisitor.symTab[v.Name] = v.Type //output streams must be declared in order to be used by other output streams
	}

	for _, v := range spec.Output {
		typeVisitor.reqType = v.Type //we mark the type that the overall expression must have, the declared type of the output stream
		v.Expr.Accept(&typeVisitor)
	}

	for _, e := range typeVisitor.errors {
		fmt.Printf(prefix+"Error %s\n", e)
	}
}
