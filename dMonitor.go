package dLola

import (
	//	"errors"
	"fmt"
	//	"strconv"
)

type Time = int
type Id = int
type MsgType int

const (
	Res     MsgType = iota
	Req     MsgType = iota
	Trigger MsgType = iota
)

//instance Show MsgType where
//  show Res = "Res"
//  show Req = "Req"
//  show Trigger = "Trigger"

func (msg MsgType) String() string {
	// declare an array of strings operator counts how many items in the array (7)
	names := [...]string{
		"Res",
		"Req",
		"Trigger"}
	if msg < Res || msg > Trigger {
		return "Unknown"
	}
	return names[msg]
}

type Msg struct {
	Kind        MsgType
	Stream      InstStreamExpr
	Value       *InstExpr //instead of LolaType, may hold any value of any type supported by golang, using type assertions the type can be checked and the value can be retrieved
	ResTime     *Time
	SimplRounds *SimplRounds
	Src         Id
	Dst         Id
}

func (msg Msg) String() string {
	val := "nil"
	if msg.Value != nil {
		v := *(msg.Value)
		val = v.Sprint()
	}
	time := "nil"
	if msg.ResTime != nil {
		t := *(msg.ResTime)
		time = fmt.Sprintf("%d", t)
	}
	simp := "nil"
	if msg.SimplRounds != nil {
		s := *(msg.SimplRounds)
		simp = fmt.Sprintf("%d", s)
	}
	return fmt.Sprintf("Msg{ kind = %s\nstream = %s\nvalue = %s\nresTime = %s\nsimplRounds = %s\nsrc = %d\ndst = %d\n}", msg.Kind.String(), msg.Stream.Sprint(), val, time, simp, msg.Src, msg.Dst)
	//return fmt.Sprintf("%v", msg)
}

func Equal(msg Msg, msg2 Msg) bool {
	return msg == msg2
}

type Received = []Msg                        //[Msg]
type Pending = []Msg                         //[Msg]
type Output = []Msg                          //[Msg]
type Requested = map[InstStreamExpr]struct{} //S.Set Stream

//type Type = Int // 0 Bool, 1 Float
type SimplRounds = int
type Eval = bool //won't need a Req to be sent to its destiny, will be flipped after creating the Msg to avoid duplicates
type Resp struct {
	value       InstExpr //LolaType
	eval        Eval
	resTime     Time
	simplRounds SimplRounds
} //result, Eval|Lazy, time at which the result was obtained, #of calls to simplExp
type Resolved struct {
	stream InstStreamExpr
	resp   Resp
}
type Und struct {
	exp         InstExpr
	eval        Eval
	simplRounds SimplRounds
}
type Unresolved struct {
	stream InstStreamExpr
	und    Und
}
type ExpEval struct {
	exp  InstExpr
	eval Eval
}

type RSet = map[InstStreamExpr]Resp //M.Map Stream Resp
type USet = map[InstStreamExpr]Und  //M.Map Stream Und
//type ExpSet = Spec                  //M.Map Stream (Exp, Eval)

type Monitor struct {
	nid            Id
	q              Received
	i              []Resolved
	u              USet
	r              RSet
	expr           Spec //contains eval Streams
	pen            Pending
	out            Output
	req            Requested
	t              Time
	routes         map[Id]Id         //M.Map Id Id
	delta          map[StreamName]Id //in which node the stream will be computed
	tracelen       int
	numMsgs        int
	sumPayload     int
	redirectedMsgs int                 // part of numMsgs
	depGraph       DepGraphAdj         //dependencies among streams
	dep            map[StreamName][]Id //Monitors that need the value of the stream (should be coherent to delta)
	trigger        []Resolved
}

func (n Monitor) String() string {
	s := fmt.Sprintf("\n###############\n Node { nid = %d\n q = %s\n i = %v\n u = %v\n r = %v\n expr = %v\n pen = %v\n out = %v\n req = %v\n t = %d\n routes = %v\n delta = %v\n tracelen = %v\n"+
		" numMsgs = %d\n sumPayload = %d\n redirectedMsgs = %d\n dep = %v\n trigger = %v\n} ################\n",
		n.nid, n.q, n.i, printU(n.u), n.r, PrettyPrintSpec(&(n.expr), ""), n.pen, n.out, n.req, n.t, n.routes, n.delta, n.tracelen, n.numMsgs, n.sumPayload, n.redirectedMsgs, n.dep, n.trigger)
	return s
}
func (m Monitor) triggered() bool {
	return len(m.trigger) > 0
}
func (m Monitor) finished() bool {
	return m.tracelen <= m.t && len(m.q) == 0 && len(m.i) == 0 && len(m.pen) == 0 && len(m.out) == 0
}
func printU(u USet) string {
	s := "map:["
	for istream, und := range u {
		s += fmt.Sprintf("%s : {%s, %t, %d}; ", istream.Sprint(), und.exp.Sprint(), und.eval, und.simplRounds)
	}
	s += "]"
	return s
}

func PrintMons(ms map[Id]*Monitor) string {
	s := ""
	for _, m := range ms {
		s += m.String()
	}
	return s
}

func NewMonitor(id, tracelen int, s Spec, received Received, routes map[Id]Id, delta map[StreamName]Id, depGraph DepGraphAdj, dep map[StreamName][]Id) Monitor {
	return Monitor{id, received, make([]Resolved, 0), make(USet), make(RSet), s, Pending{}, Output{}, Requested{}, 0, routes, delta, tracelen, 0, 0, 0, depGraph, dep, make([]Resolved, 0)}
}

func Converge(mons map[Id]*Monitor) {
	allfinished := false
	anytriggered := false
	fmt.Printf("Converge: %s\n", PrintMons(mons))
	for !(anytriggered || allfinished) {
		fmt.Printf("Should continue converging\n")
		Tick(mons)
		anytriggered, allfinished = ShouldContinue(mons)
		//fmt.Printf("Should continue converging because %t, %t\n", anytriggered, allfinished)
	}
	fmt.Printf("Finished\n")
	return

}

func ShouldContinue(mons map[Id]*Monitor) (bool, bool) {
	allfinished := true
	anytriggered := false
	for _, m := range mons {
		if m.triggered() {
			anytriggered = true
			break
		}
		allfinished = allfinished && m.finished()
	}
	return anytriggered, allfinished
}

func Tickn(mons map[Id]*Monitor, nticks int) {
	for i := 0; i < nticks; i++ {
		Tick(mons)
	}
}

func Tick(mons map[Id]*Monitor) {
	fmt.Printf("Tick mons:%s\n", PrintMons(mons))
	for _, m := range mons {
		m.process()
	}
	for _, m := range mons {
		//fmt.Printf("Before dispatch of mon %d:%s\n", m.nid, PrintMons(mons))
		m.dispatch(mons)
		//fmt.Printf("After dispatching of mon %d:%s\n", m.nid, PrintMons(mons))
	}
	fmt.Printf("TICKED mons:%s\n", PrintMons(mons))
}

func (m *Monitor) dispatch(mons map[Id]*Monitor) {
	for _, msg := range m.out {
		nextHopMon := m.routes[msg.Dst]                      //we look for the nextHop in the route from m to msg.Dst
		mons[nextHopMon].q = append(mons[nextHopMon].q, msg) //append msg to the incoming msgs of the destination
	}
	m.out = []Msg{}
}

func (m *Monitor) process() {
	m.processQ()
	m.readInput()
	m.generateEquations()
	m.simplify()
	m.addRes()
	m.addReq()
	m.t++
}

func (m *Monitor) processQ() {
	//fmt.Printf("[%d]:PROCESSQ: %s\n", m.nid, m.String())
	for _, msg := range m.q {
		if msg.Dst != m.nid { //redirect msgs whose dst is not this node
			m.out = append(m.out, msg)
			m.numMsgs++
			m.sumPayload += payload(msg)
		} else {
			switch msg.Kind {
			case Res:
				m.r[msg.Stream] = msgToResp(msg)
			case Req:
				m.pen = append(m.pen, msg)
			case Trigger:
				m.pen = append(m.pen, msg)

			}
		}
	}
	m.q = []Msg{}
}

func msgToResp(msg Msg) Resp {
	return Resp{*msg.Value, false, *msg.ResTime, *msg.SimplRounds} //received responses are marked as LAZY, so as not to send them again and flood the net!!
}

func (m *Monitor) readInput() {
	//fmt.Printf("[%d]:READINPUT: %s\n", m.nid, m.String())
	//TODO: use channels one for each input, which will provide the input streams events
	//those events may be fed by some generator, a file, a socket... running in a go routine
}
func (m *Monitor) generateEquations() {
	//fmt.Printf("[%d]:GENERATEEQ: %s\n", m.nid, m.String())
	if m.t < m.tracelen { //expressions will be intantiated from tick 0 to tracelen - 1
		for _, o := range m.expr.Output {
			if m.delta[o.Name] == m.nid {
				e := o.Expr
				i := e.InstantiateExpr(m.t, m.tracelen)
				//fmt.Printf("Instanced expr: %s\n", i.Sprint())
				u := Und{SimplifyExpr(i), o.Eval, 0} //Und gets generated with the eval specified
				stream := InstStreamFetchExpr{o.Name, m.t}
				m.u[stream] = u
			}
		}
	}
}
func (m *Monitor) simplify() {
	//fmt.Printf("[%d]:SIMPLIFY: %s\n", m.nid, m.String())
	someSimpl := true //will control when to stop searching for substitutions and
	for someSimpl {
		//fmt.Printf("Again\n")
		someSimpl = false
		for stream, und := range m.u { //for each unresolved stream in U
			dep := m.depGraph[stream.GetName()]
			for _, d := range dep { //for each of its dependencies
				depStream := InstStreamFetchExpr{StreamName(d.Dest), stream.GetTick() + d.Weight} //build the depStream taking into account the tick of the stream and the weight of the dependency
				//fmt.Printf("need the stream %s to simplify\n", depStream.Sprint())
				resp, ok := m.r[depStream]
				if ok { //we found the value of the dependency stream in R
					//fmt.Printf("and we have its value\n")
					und.exp = und.exp.Substitute(depStream, resp.value)
				}
			}
			und.exp = SimplifyExpr(und.exp)
			und.simplRounds++
			newresp, isResolved := toResp(m.t, und)
			if isResolved {
				//fmt.Printf("New Resp: %v", newresp)
				m.r[stream] = newresp //add it to R
				delete(m.u, stream)   //remove it from U
			}
			someSimpl = someSimpl || isResolved
		}
	}
}

func toResp(t int, und Und) (Resp, bool) {
	//fmt.Printf("To Resp: %v\n", und.exp)
	var r Resp
	switch und.exp.(type) {
	case InstTruePredicate:
		return Resp{und.exp, und.eval, t, und.simplRounds}, true
	case InstFalsePredicate:
		return Resp{und.exp, und.eval, t, und.simplRounds}, true
	case InstIntLiteralExpr:
		return Resp{und.exp, und.eval, t, und.simplRounds}, true
	case InstFloatLiteralExpr:
		return Resp{und.exp, und.eval, t, und.simplRounds}, true
	case InstStringLiteralExpr:
		return Resp{und.exp, und.eval, t, und.simplRounds}, true
	default:
		//fmt.Printf("Not resolved")
		return r, false
	}
}

func (m *Monitor) addRes() {
	//fmt.Printf("[%d]:ADDRES: %s\n", m.nid, m.String())
	for stream, resp := range m.r {
		if resp.eval && stream.GetTick() <= m.t && stream.GetTick() < m.tracelen { //EVAL streams
			if destinies, ok := m.dep[stream.GetName()]; ok {
				for _, d := range destinies {
					if d != m.nid {
						msg := createMsg(stream, &resp, m.nid, d)
						fmt.Printf("Creating Res msg of eval stream %s\n", msg.String())
						m.out = append(m.out, msg)
					}
				}
				resp.eval = false //we mark the resp as already sent to interested monitors
			}
		}
	}
	newPen := make([]Msg, 0)
	for _, msg := range m.pen { //LAZY streams need to be requested in order to send responses
		if resp, ok := m.r[msg.Stream]; ok { //note this msg will no longer be in pen
			fmt.Printf("Resolved LAZY %s\n", msg.String())
			newMsg := createMsg(msg.Stream, &resp, m.nid, msg.Src)
			m.out = append(m.out, newMsg)
			if msg.Kind == Trigger {
				fmt.Printf("Resolved Trigger %s\n", msg.Stream.Sprint())
				m.trigger = append(m.trigger, Resolved{msg.Stream, resp})
			}
		} else {
			newPen = append(newPen, msg) //if we do not have the resp, we will keep it in pen
		}
	}
	m.pen = newPen
}

func createMsg(stream InstStreamExpr, resp *Resp, id, dst Id) Msg {
	if resp == nil {
		return Msg{Req, stream, nil, nil, nil, id, dst}
	}
	return Msg{Res, stream, &resp.value, &resp.resTime, &resp.simplRounds, id, dst}
}
func (m *Monitor) addReq() {
	fmt.Printf("[%d]:ADDREQ: %s\n", m.nid, m.String())
	for _, p := range m.pen {
		createReqMsgsPen(p.Stream, m)
	}
}

func createReqMsgsPen(stream InstStreamExpr, m *Monitor) {
	fmt.Printf("Creating REQS\n")
	adjacencies := m.depGraph[stream.GetName()]
	dependencies := convertToStreams(stream, adjacencies, m.tracelen)
	for i := 0; i < len(dependencies); i++ {
		depStream := dependencies[i]
		//fmt.Printf("Adjacency %s\n", adj.Sprint())
		createReqStream(depStream, m)
		dependencies = addNextLevelDependencies(dependencies, convertToStreams(depStream, m.depGraph[depStream.GetName()], m.tracelen))
		fmt.Printf("next level dependencies after: %s\n", SprintStreams(dependencies))
	}
}

func convertToStreams(stream InstStreamExpr, adjacencies []Adj, tlen int) []InstStreamExpr {
	r := make([]InstStreamExpr, 0)
	for _, adj := range adjacencies {
		depTick := stream.GetTick() + adj.Weight
		if depTick >= 0 && depTick < tlen { //the dependency could be instantiated
			depStream := InstStreamFetchExpr{adj.Dest, depTick}
			r = append(r, depStream)
		}
	}
	return r
}

/*will create a req msg for depStream and add it to out iff the stream is not in R, assigned to other monitor(delta), not previously requested, not eval and have been instanced*/
func createReqStream(depStream InstStreamExpr, m *Monitor) {
	resp, resolved := m.r[depStream]
	_, requested := m.req[depStream]
	fmt.Printf("Dependency could be intantiated tlen: %d\n%s\n !resolved %t, !requested %t, !eval %t, assigned to other monitor: %t\n", m.tracelen, depStream.Sprint(), !resolved, !requested, !resp.eval, m.delta[depStream.GetName()] != m.nid)
	if !resolved && m.delta[depStream.GetName()] != m.nid && !resp.eval && !requested && depStream.GetTick() <= m.t { //not in R, not assigned to this monitor, not eval and not already requested, allow request of streams that have not yet been instanced?
		fmt.Printf("Creatting Request: %s\n", depStream.Sprint())
		m.out = append(m.out, createMsg(depStream, nil, m.nid, m.delta[depStream.GetName()]))
		m.req[depStream] = struct{}{} //mark it as requested
	}
}

/*will change dependencies*/
func addNextLevelDependencies(dependencies, candidates []InstStreamExpr) []InstStreamExpr {
	fmt.Printf("next level dependencies before: %s\ncandidates: %s\n", SprintStreams(dependencies), SprintStreams(candidates))
	for _, c := range candidates {
		if !elemStream(dependencies, c, EqInstStreamExpr) { //add if not already present
			dependencies = append(dependencies, c)
		}
	}
	fmt.Printf("next level dependencies after: %s\n", SprintStreams(dependencies))
	return dependencies
}

func payload(msg Msg) int {
	payload := 0
	val := *msg.Value
	switch v := val.(type) {
	case InstTruePredicate:
		payload = commonPayLoad(msg.Stream.GetName().Sprint()) + 1
	case InstFalsePredicate:
		payload = commonPayLoad(msg.Stream.GetName().Sprint()) + 1
	case InstIntLiteralExpr:
		payload = commonPayLoad(msg.Stream.GetName().Sprint()) + 32
	case InstFloatLiteralExpr:
		payload = commonPayLoad(msg.Stream.GetName().Sprint()) + 32
	case InstStringLiteralExpr:
		payload = commonPayLoad(msg.Stream.GetName().Sprint()) + len(v.S)*8
	default:

	}
	return payload
}

func commonPayLoad(s string) int {
	return 2 + 32*5 + 8*len(s)
}

/*in bits
payLoad :: Int -> Msg -> Int
payLoad acc m = case (stream m, value m) of
  ((s, t), (Bt x)) -> acc + commonPayLoad s 1
  ((s, t), (Bt3 x)) -> acc + commonPayLoad s 2
  ((s, t), (Bt4 x)) -> acc + commonPayLoad s 2
  ((s, t), (Nt x)) -> acc + commonPayLoad s basicTypeSize
  ((s, t), (Mt x)) -> acc + commonPayLoad s (basicTypeSize * M.size x)
  ((s, t), (St x)) -> acc + commonPayLoad s (basicTypeSize * S.size x)

basicTypeSize = 64
--in bits
{-
1 3-valued : kind
3 int : time, src and dst ### resTime and simplRounds will not be counted because they are here fro profiling purposes only
1 string of characters of 8 bits
1 typevalue of arbitrary size
-}
commonPayLoad :: String -> Int -> Int
commonPayLoad string valuePayLoad = 2 + 32 *5 + 8 * length string + valuePayLoad
*/
