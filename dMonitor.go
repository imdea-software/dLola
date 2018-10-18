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
	Value       InstExpr //instead of LolaType, may hold any value of any type supported by golang, using type assertions the type can be checked and the value can be retrieved
	ResTime     Time
	SimplRounds SimplRounds
	Src         Id
	Dst         Id
}

func (msg Msg) String() string {
	/*s := "Msg{ kind = " + msg.Kind.String() + "\n" +
		"stream = " + msg.Stream + "\n" +
		"value = " + msg.Value + "\n" +
		"resTime = " + msg.ResTime + "\n" +
		"simplRounds = " + msg.SimplRounds + "\n" +
		"src = " + msg.Src + "\n" +
		"dst = " + msg.Dst + "\n}"
	return s*/
	return fmt.Sprintf("%v", msg)
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
	expr           Spec
	pen            Pending
	out            Output
	req            Requested
	t              Time
	routes         map[Id]Id         //M.Map Id Id
	delta          map[StreamName]Id //in which node the stream will be computed
	tracelen       int
	numMsgs        int
	sumPayload     int
	redirectedMsgs int // part of numMsgs
	evalStreams    []StreamName
	depGraph       DepGraphAdj         //dependencies among streams
	dep            map[StreamName][]Id //Monitors that need the value of the stream (should be coherent to delta)
	trigger        []Resolved
}

func (n Monitor) string() string {
	s := fmt.Sprintf("\n###############\n Node { nid = %d\n q = %s\n i = %v\n u = %v\n r = %v\n expr = %v\n pen = %v\n out = %v\n req = %v\n t = %d\n routes = %v\n delta = %v\n tracelen = %v\n"+
		" numMsgs = %d\n sumPayload = %d\n redirectedMsgs = %d\n evalStreams = %v\n dep = %v\n trigger = %v\n} ################\n",
		n.nid, n.q, n.i, printU(n.u), n.r, PrettyPrintSpec(&(n.expr), ""), n.pen, n.out, n.req, n.t, n.routes, n.delta, n.tracelen, n.numMsgs, n.sumPayload, n.redirectedMsgs, n.evalStreams, n.dep, n.trigger)
	return s
}

func printU(u USet) string {
	s := ""
	for istream, und := range u {
		s += fmt.Sprintf("%s : {%s, %t, %d}; ", istream.Sprint(), und.exp.Sprint(), und.eval, und.simplRounds)
	}
	return s
}

func PrintMons(ms map[Id]*Monitor) string {
	s := ""
	for _, m := range ms {
		s += m.string()
	}
	return s
}
func (m Monitor) triggered() bool {
	return len(m.trigger) > 0
}
func (m Monitor) finished() bool {
	return m.tracelen <= m.t && len(m.q) == 0 && len(m.i) == 0 && len(m.pen) == 0 && len(m.out) == 0
}

func NewMonitor(id, tracelen int, s Spec, routes map[Id]Id, delta map[StreamName]Id, eval []StreamName, depGraph DepGraphAdj, dep map[StreamName][]Id) Monitor {
	return Monitor{id, Received{}, make([]Resolved, 0), make(USet), make(RSet), s, Pending{}, Output{}, Requested{}, 0, routes, delta, tracelen, 0, 0, 0, eval, depGraph, dep, make([]Resolved, 0)}
}

func RoundrDelta(s Spec, nmons int) map[StreamName]Id {
	delta := make(map[StreamName]Id)
	var i Id = 0
	for _, o := range s.Output {
		delta[o.Name] = i
		i = (i + 1) % nmons
	}
	return delta
}

func InterestedMonitors(delta map[StreamName]Id, depGraph DepGraphAdj) map[StreamName][]Id {
	monitorDependencies := make(map[StreamName][]Id)
	for stream, mon := range delta {
		streamDeps := depGraph[stream]
		for _, streamDep := range streamDeps {
			monitorDependencies[streamDep.Dest] = append(monitorDependencies[streamDep.Dest], mon) //the Monitor assigned to compute the stream stream will need the value of streamDep.Dst TODO revise append to nil
		}
	}
	return monitorDependencies
}

func Converge(mons map[Id]*Monitor) {
	allfinished := false
	anytriggered := false
	for !(anytriggered || allfinished) {
		Tick(mons)
		anytriggered, allfinished = ShouldContinue(mons)
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
	//fmt.Printf("Converge mons:%s\n", printMons(mons))
	for _, m := range mons {
		m.process()
	}
	for _, m := range mons {
		m.dispatch(mons)
	}
	//fmt.Printf("TICKED mons:%s\n", printMons(mons))
	//tickNode in process()
}

func (m *Monitor) dispatch(mons map[Id]*Monitor) {
	for _, msg := range m.out {
		nextHopMon := m.routes[msg.Dst] //we look for the nextHop in the route from m to msg.Dst
		destq := mons[nextHopMon].q
		destq = append(destq, msg) //append msg to the incoming msgs of the destination
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
	for _, msg := range m.q {
		if msg.Dst != m.nid { //redirect msgs whose dst is not this node
			m.out = append(m.out, msg)
			m.numMsgs++
			m.sumPayload += payload(msg)
		} else {
			if msg.Kind == Res {
				m.r[msg.Stream] = msgToResp(msg)
			} else {
				if msg.Kind == Req {
					m.pen = append(m.pen, msg)
				}
			}
		}
	}

}

func msgToResp(msg Msg) Resp {
	return Resp{msg.Value, false, msg.ResTime, msg.SimplRounds} //received responses are marked as LAZY, so as not to send them again and flood the net!!
}

func (m *Monitor) readInput() {
	//TODO: use channels one for each input, which will provide the input streams events
	//those events may be fed by some generator, a file, a socket... running in a go routine
}
func (m *Monitor) generateEquations() {
	eval := false //eval should be specified in spec
	if m.t <= m.tracelen {
		for _, o := range m.expr.Output {
			e := o.Expr
			i := e.InstantiateExpr(m.t, m.tracelen)
			u := Und{SimplifyExpr(i), eval, 0}
			stream := InstStreamFetchExpr{o.Name, m.t}
			m.u[stream] = u
		}
	}
}
func (m *Monitor) simplify() {
	someSimpl := true //will control when to stop searching for substitutions and
	for someSimpl {
		//fmt.Printf("Again\n")
		someSimpl = false
		for stream, und := range m.u { //for each unresolved stream in U
			dep := m.depGraph[stream.GetName()]
			for _, d := range dep { //for each of its dependencies
				depStream := InstStreamFetchExpr{StreamName(d.Dest), m.t + d.Weight}
				resp, ok := m.r[depStream]
				if ok { //we found the value of the dependency stream in R
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
	for stream, resp := range m.r {
		if resp.eval && stream.GetTick() <= m.t && stream.GetTick() < m.tracelen { //EVAL streams
			if destinies, ok := m.dep[stream.GetName()]; ok {
				for _, d := range destinies {
					if d != m.nid {
						msg := createMsg(stream, &resp, m.nid, d)
						m.out = append(m.out, msg)
					}
				}
				resp.eval = false //we mark the resp as already sent to interested monitors
			}
		}
	}
	for _, msg := range m.pen { //LAZY streams
		if resp, ok := m.r[msg.Stream]; ok {
			msg := createMsg(msg.Stream, &resp, m.nid, msg.Src)
			m.out = append(m.out, msg)
			if msg.Kind == Trigger {
				m.trigger = append(m.trigger, Resolved{msg.Stream, resp})
			}
		}
	}

}

func createMsg(stream InstStreamExpr, resp *Resp, id, dst Id) Msg {
	if resp == nil {
		var value InstExpr
		var time Time
		var simplRounds int
		return Msg{Req, stream, value, time, simplRounds, id, dst}
	}
	return Msg{Res, stream, resp.value, resp.resTime, resp.simplRounds, id, dst}
}
func (m *Monitor) addReq() {
	for _, p := range m.pen {
		if _, ok := m.u[p.Stream]; ok { //pending stream is in U
			dependencies := m.depGraph[p.Stream.GetName()]
			convertFilterResolvedLocal(p.Stream, dependencies, m)
			//			r = append(r, filtered...)
		}
	}
}

/*func getNeededStreams(m *Monitor) []InstStreamExpr {
	r := make([]InstStreamExpr, 0)
	return r
}*/

func convertFilterResolvedLocal(stream InstStreamExpr, dependencies []Adj, m *Monitor) {
	for _, adj := range dependencies {
		depTick := stream.GetTick() + adj.Weight
		if depTick >= 0 && depTick < m.tracelen { //the dependency could be instantiated
			depStream := InstStreamFetchExpr{adj.Dest, depTick}
			resp, resolved := m.r[depStream]
			_, requested := m.req[depStream]
			if !resolved && m.delta[adj.Dest] != m.nid && !resp.eval && !requested { //not in R, not assigned to this monitor, not eval and not already requested
				var r Resp
				m.out = append(m.out, createMsg(depStream, &r, m.nid, m.delta[adj.Dest]))
				m.req[depStream] = struct{}{}
			}
		}
	}
}

func payload(msg Msg) int {
	payload := 0
	switch v := msg.Value.(type) {
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
